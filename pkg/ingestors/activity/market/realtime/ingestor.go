package realtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

type Ingestor struct {
	store            market.MarketData
	exchangeRegistry registry.ConnectorRegistry
	assetRegistry    registry.AssetRegistry
	logger           logging.ApplicationLogger
	healthStore      health.HealthStore

	// WebSocket management
	wsContext context.Context
	wsCancel  context.CancelFunc
	isActive  bool
	mutex     sync.RWMutex

	activeConnections map[connector.ExchangeName]connector.WebSocketConnector

	// Track which instrument types are subscribed per asset
	subscriptions     map[portfolio.Asset][]connector.Instrument
	subscriptionMutex sync.RWMutex
}

func NewIngestor(
	store market.MarketData,
	exchangeRegistry registry.ConnectorRegistry,
	assetRegistry registry.AssetRegistry,
	logger logging.ApplicationLogger,
	healthStore health.HealthStore,
) *Ingestor {
	return &Ingestor{
		store:             store,
		exchangeRegistry:  exchangeRegistry,
		assetRegistry:     assetRegistry,
		logger:            logger,
		healthStore:       healthStore,
		activeConnections: make(map[connector.ExchangeName]connector.WebSocketConnector),
		subscriptions:     make(map[portfolio.Asset][]connector.Instrument),
	}
}

func (ri *Ingestor) Start(ctx context.Context) error {
	ri.mutex.Lock()
	defer ri.mutex.Unlock()

	if ri.isActive {
		return fmt.Errorf("realtime ingestor already active")
	}

	ri.wsContext, ri.wsCancel = context.WithCancel(ctx)
	ri.isActive = true

	// Get required assets from strategy configs
	tradingAssets := ri.assetRegistry.GetRequiredAssets()
	ri.logger.Info("📋 Required trading assets from strategies: %d assets", len(tradingAssets))
	for _, asset := range tradingAssets {
		ri.logger.Info("  - %s", asset.Symbol())
	}

	connectors := ri.exchangeRegistry.GetTradingWebSocketConnectors()

	if len(tradingAssets) == 0 {
		ri.logger.Info("❌ No assets required by enabled strategies - realtime ingestor won't start")
		return nil
	}

	ri.logger.Info("🔄 About to start WebSocket streams for trading assets...")

	// Start WebSocket streams for each exchange
	ri.logger.Info("🔗 Found %d trading WebSocket connectors", len(connectors))

	for _, conn := range connectors {
		exchangeName := conn.GetConnectorInfo().Name
		ri.logger.Info("🚀 Starting WebSocket stream for %s", exchangeName)
		go ri.startExchangeStream(conn, tradingAssets)
	}

	ri.logger.Info("✅ Started WebSocket ingestion for %d assets", len(tradingAssets))
	return nil
}

func (ri *Ingestor) startExchangeStream(wsConn connector.WebSocketConnector, assets []portfolio.Asset) {
	exchangeName := wsConn.GetConnectorInfo().Name

	// Start WebSocket connection
	if err := wsConn.StartWebSocket(ri.wsContext); err != nil {
		ri.logger.Error("Failed to start WebSocket for %s: %v", exchangeName, err)
		return
	}

	// Get detailed asset requirements (which instrument types per asset)
	assetRequirements := ri.assetRegistry.GetAssetRequirements()

	// Subscribe to data based on actual requirements
	for _, req := range assetRequirements {
		// Subscribe to orderbooks for each required instrument type
		for _, instrumentType := range req.Instruments {
			if err := wsConn.SubscribeOrderBook(req.Asset, instrumentType); err != nil {
				ri.logger.Error("Failed to subscribe to %s orderbook for %s on %s: %v",
					instrumentType, req.Asset.Symbol(), exchangeName, err)
			} else {
				ri.logger.Info("✅ Subscribed to %s orderbook for %s on %s",
					instrumentType, req.Asset.Symbol(), exchangeName)

				// Track this subscription
				ri.subscriptionMutex.Lock()
				ri.subscriptions[req.Asset] = append(ri.subscriptions[req.Asset], instrumentType)
				ri.subscriptionMutex.Unlock()
			}
		}

		// Subscribe to klines for strategy analysis
		klineIntervals := []string{"1m", "5m", "15m", "1h"}
		for _, interval := range klineIntervals {
			if err := wsConn.SubscribeKlines(req.Asset, interval); err != nil {
				ri.logger.Error("Failed to subscribe to %s klines for %s on %s: %v",
					interval, req.Asset.Symbol(), exchangeName, err)
			}
		}
	}

	// Process real-time updates
	go ri.processOrderBookStream(wsConn, exchangeName)
	go ri.processKlineStream(wsConn, exchangeName)
	go ri.processErrorStream(wsConn, exchangeName)
}

func (ri *Ingestor) processKlineStream(wsConn connector.WebSocketConnector, exchangeName connector.ExchangeName) {
	ri.logger.Info("🔄 Starting kline stream processing for %s", exchangeName)

	klineChan := wsConn.KlineUpdates()
	ri.logger.Info("📊 Got kline channel for %s, waiting for updates...", exchangeName)

	for {
		select {
		case <-ri.wsContext.Done():
			ri.logger.Info("❌ Context cancelled, stopping kline stream for %s", exchangeName)
			return
		case klineUpdate, ok := <-klineChan:
			if !ok {
				ri.logger.Info("📪 Kline channel closed for %s", exchangeName)
				return
			}

			asset := portfolio.NewAsset(klineUpdate.Symbol)
			ri.logger.Info("📊 Received %s kline update for %s on %s", klineUpdate.Interval, klineUpdate.Symbol, exchangeName)
			ri.logger.Info("💾 Storing kline in asset store: %s/%s - O:%.2f H:%.2f L:%.2f C:%.2f",
				exchangeName, klineUpdate.Symbol,
				klineUpdate.Open.InexactFloat64(), klineUpdate.High.InexactFloat64(),
				klineUpdate.Low.InexactFloat64(), klineUpdate.Close.InexactFloat64())

			ri.store.UpdateKline(asset, exchangeName, klineUpdate)

			// Report successful data receipt to health monitoring
			ri.healthStore.RecordDataReceived(exchangeName, health.DataTypeKlines, health.SourceWebSocket, 0)

			// CRITICAL: Update market data when klines arrive to refresh orderbook/prices
			if exchange, exists := ri.exchangeRegistry.GetConnector(exchangeName); exists {
				// Access the market simulator directly through interface with proper method signature
				if marketUpdater, ok := exchange.(interface{ UpdateMarketData(time.Time) error }); ok {
					if err := marketUpdater.UpdateMarketData(klineUpdate.CloseTime); err != nil {
						ri.logger.Error("Failed to update market data for %s: %v", exchangeName, err)
					} else {
						ri.logger.Info("🔄 Updated market data for %s with kline timestamp %v", exchangeName, klineUpdate.CloseTime)
					}
				} else {
					ri.logger.Error("❌ Exchange %s does not implement UpdateMarketData method", exchangeName)
				}
			} else {
				ri.logger.Error("❌ Exchange connector not found for %s", exchangeName)
			}
		}
	}
}

func (ri *Ingestor) processOrderBookStream(wsConn connector.WebSocketConnector, exchangeName connector.ExchangeName) {
	ri.logger.Info("🔄 Starting orderbook stream processing for %s", exchangeName)

	orderBookChan := wsConn.OrderBookUpdates()
	ri.logger.Info("📡 Got orderbook channel for %s, waiting for updates...", exchangeName)

	for {
		select {
		case <-ri.wsContext.Done():
			ri.logger.Info("❌ Context cancelled, stopping orderbook stream for %s", exchangeName)
			return
		case orderBookUpdate, ok := <-orderBookChan:
			if !ok {
				ri.logger.Info("📪 OrderBook channel closed for %s", exchangeName)
				return
			}

			// Get which instrument types we subscribed to for this asset
			ri.subscriptionMutex.RLock()
			instrumentTypes := ri.subscriptions[orderBookUpdate.Asset]
			ri.subscriptionMutex.RUnlock()

			// Store the orderbook update for each subscribed instrument type
			for _, instrumentType := range instrumentTypes {
				ri.store.UpdateOrderBook(
					orderBookUpdate.Asset,
					exchangeName,
					instrumentType,
					orderBookUpdate,
				)

				ri.logger.Debug("📊 Updated %s orderbook for %s on %s",
					instrumentType, orderBookUpdate.Asset.Symbol(), exchangeName)
			}

			// Report successful data receipt to health monitoring
			ri.healthStore.RecordDataReceived(exchangeName, health.DataTypeOrderbooks, health.SourceWebSocket, 0)
		}
	}
}

func (ri *Ingestor) processErrorStream(wsConn connector.WebSocketConnector, exchangeName connector.ExchangeName) {
	for {
		select {
		case <-ri.wsContext.Done():
			return
		case err := <-wsConn.ErrorChannel():
			ri.logger.Error("WebSocket error for %s: %v", exchangeName, err)
			// Report error to health monitoring - affects all data types on this websocket
			ri.healthStore.RecordDataError(exchangeName, health.DataTypeKlines, err)
			ri.healthStore.RecordDataError(exchangeName, health.DataTypeOrderbooks, err)
		}
	}
}

func (ri *Ingestor) Stop() error {
	ri.mutex.Lock()
	defer ri.mutex.Unlock()

	if !ri.isActive {
		return nil
	}

	if ri.wsCancel != nil {
		ri.wsCancel()
	}

	// Stop all WebSocket connections
	for _, conn := range ri.exchangeRegistry.GetTradingWebSocketConnectors() {
		conn.StopWebSocket()
	}

	ri.isActive = false
	ri.logger.Info("Stopped WebSocket ingestion")
	return nil
}

func (ri *Ingestor) IsActive() bool {
	ri.mutex.RLock()
	defer ri.mutex.RUnlock()
	return ri.isActive
}

func (ri *Ingestor) GetActiveConnections() map[connector.ExchangeName]connector.WebSocketConnector {
	ri.mutex.RLock()
	defer ri.mutex.RUnlock()

	result := make(map[connector.ExchangeName]connector.WebSocketConnector)
	for name, conn := range ri.activeConnections {
		result[name] = conn
	}
	return result
}
