package ingestors

import (
	"context"
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// RealtimeIngestor is a generic base implementation for WebSocket data collection
type RealtimeIngestor struct {
	conn          interface{}
	wsCapable     connector.WebSocketCapable
	wsSubscriber  realtime.WebSocketSubscriber
	exchangeName  connector.ExchangeName
	marketType    connector.MarketType
	assetRegistry registry.AssetRegistry
	store         market.MarketStore
	logger        logging.ApplicationLogger

	// State
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
	mu       sync.RWMutex

	// Extension point for market-specific WebSocket subscriptions
	extensions []realtime.WebSocketExtension
}

func NewRealtimeIngestor(
	conn interface{},
	exchangeName connector.ExchangeName,
	marketType connector.MarketType,
	assetRegistry registry.AssetRegistry,
	store interface{},
	logger logging.ApplicationLogger,
	extensions ...realtime.WebSocketExtension,
) realtime.RealtimeIngestor {
	// Cast to WebSocketCapable for lifecycle management
	wsCapable, ok := conn.(connector.WebSocketCapable)
	if !ok {
		logger.Error("Connector does not implement WebSocketCapable interface")
		return nil
	}

	// Cast to WebSocketSubscriber for subscriptions
	wsSubscriber, ok := conn.(realtime.WebSocketSubscriber)
	if !ok {
		logger.Error("Connector does not implement WebSocketSubscriber interface")
		return nil
	}

	// Cast store to MarketStore for writing data
	marketStore, ok := store.(market.MarketStore)
	if !ok {
		logger.Error("Store does not implement MarketStore interface")
		return nil
	}

	return &RealtimeIngestor{
		conn:          conn,
		wsCapable:     wsCapable,
		wsSubscriber:  wsSubscriber,
		exchangeName:  exchangeName,
		marketType:    marketType,
		assetRegistry: assetRegistry,
		store:         marketStore,
		logger:        logger,
		extensions:    extensions,
	}
}

func (ri *RealtimeIngestor) Start(ctx context.Context) error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if ri.isActive {
		return fmt.Errorf("realtime ingestor for %s already active", ri.exchangeName)
	}

	assets := ri.assetRegistry.GetRequiredAssets()
	if len(assets) == 0 {
		ri.logger.Warn("No assets registered for %s realtime ingestion", ri.exchangeName)
		return nil
	}

	ri.ctx, ri.cancel = context.WithCancel(ctx)

	// Start WebSocket connection
	if err := ri.wsCapable.StartWebSocket(); err != nil {
		return fmt.Errorf("failed to start WebSocket for %s: %w", ri.exchangeName, err)
	}

	// Subscribe to order books for all assets
	for _, asset := range assets {
		if err := ri.wsSubscriber.SubscribeOrderBook(asset); err != nil {
			ri.logger.Error("Failed to subscribe to order book for %s on %s: %v",
				asset.Symbol(), ri.exchangeName, err)
		} else {
			ri.logger.Info("Subscribed to order book for %s on %s", asset.Symbol(), ri.exchangeName)
		}
	}

	// Subscribe to klines for all assets
	klineIntervals := []string{"1m", "5m", "15m", "1h"}
	for _, asset := range assets {
		for _, interval := range klineIntervals {
			if err := ri.wsSubscriber.SubscribeKlines(asset, interval); err != nil {
				ri.logger.Debug("Failed to subscribe to %s klines for %s on %s: %v",
					interval, asset.Symbol(), ri.exchangeName, err)
			}
		}
	}

	// Market-specific WebSocket subscriptions (funding rate updates, etc.)
	for _, ext := range ri.extensions {
		if err := ext.Subscribe(ri.conn, ri.exchangeName, assets); err != nil {
			ri.logger.Error("Failed to subscribe to %s-specific data for %s: %v",
				ri.marketType, ri.exchangeName, err)
		}
	}

	// Start processing channels
	go ri.processOrderBookChannels()
	go ri.processKlineChannels()

	// Start extension channel processing
	for _, ext := range ri.extensions {
		go ext.ProcessChannels(ri.conn, ri.exchangeName, ri.ctx)
	}

	ri.isActive = true
	ri.logger.Info("Started %s realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *RealtimeIngestor) processOrderBookChannels() {
	channels := ri.wsSubscriber.GetOrderBookChannels()
	ri.logger.Info("Processing %d order book channels for %s", len(channels), ri.exchangeName)

	var wg sync.WaitGroup
	for channelKey, orderBookChan := range channels {
		wg.Add(1)
		go func(key string, ch <-chan connector.OrderBook) {
			defer wg.Done()
			ri.processOrderBookChannel(key, ch)
		}(channelKey, orderBookChan)
	}

	wg.Wait()
	ri.logger.Info("All order book channels closed for %s", ri.exchangeName)
}

func (ri *RealtimeIngestor) processOrderBookChannel(channelKey string, orderBookChan <-chan connector.OrderBook) {
	ri.logger.Debug("Starting order book channel processor for %s on %s", channelKey, ri.exchangeName)

	for {
		select {
		case <-ri.ctx.Done():
			ri.logger.Debug("Context cancelled, stopping order book channel %s", channelKey)
			return

		case update, ok := <-orderBookChan:
			if !ok {
				ri.logger.Debug("Order book channel %s closed", channelKey)
				return
			}

			// Write to store
			ri.store.UpdateOrderBook(update.Asset, ri.exchangeName, update)
			ri.store.UpdateLastUpdated(market.UpdateKey{
				DataType: market.DataKeyOrderBooks,
				Asset:    update.Asset,
				Exchange: ri.exchangeName,
			})

			if len(update.Bids) > 0 && len(update.Asks) > 0 {
				ri.logger.Debug(
					"WebSocket updated order book for %s on %s - bid: %s, ask: %s",
					update.Asset.Symbol(),
					ri.exchangeName,
					update.Bids[0].Price.StringFixed(2),
					update.Asks[0].Price.StringFixed(2),
				)
			}
		}
	}
}

func (ri *RealtimeIngestor) processKlineChannels() {
	channels := ri.wsSubscriber.GetKlineChannels()
	ri.logger.Info("Processing %d kline channels for %s", len(channels), ri.exchangeName)

	var wg sync.WaitGroup
	for channelKey, klineChan := range channels {
		wg.Add(1)
		go func(key string, ch <-chan connector.Kline) {
			defer wg.Done()
			ri.processKlineChannel(key, ch)
		}(channelKey, klineChan)
	}

	wg.Wait()
	ri.logger.Info("All kline channels closed for %s", ri.exchangeName)
}

func (ri *RealtimeIngestor) processKlineChannel(channelKey string, klineChan <-chan connector.Kline) {
	ri.logger.Debug("Starting kline channel processor for %s on %s", channelKey, ri.exchangeName)

	for {
		select {
		case <-ri.ctx.Done():
			ri.logger.Debug("Context cancelled, stopping kline channel %s", channelKey)
			return

		case kline, ok := <-klineChan:
			if !ok {
				ri.logger.Debug("Kline channel %s closed", channelKey)
				return
			}

			asset := portfolio.NewAsset(kline.Symbol)

			// Write to store
			ri.store.UpdateKline(asset, ri.exchangeName, kline)

			ri.logger.Debug("WebSocket updated %s kline for %s on %s",
				kline.Interval, asset.Symbol(), ri.exchangeName)
		}
	}
}

func (ri *RealtimeIngestor) Stop() error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if !ri.isActive {
		return nil
	}

	// Cancel context to stop all goroutines
	if ri.cancel != nil {
		ri.cancel()
	}

	// Unsubscribe from market-specific extensions
	for _, ext := range ri.extensions {
		if err := ext.Unsubscribe(ri.conn, ri.exchangeName); err != nil {
			ri.logger.Warn("Failed to unsubscribe from %s-specific data for %s: %v",
				ri.marketType, ri.exchangeName, err)
		}
	}

	// Stop WebSocket connection
	if err := ri.wsCapable.StopWebSocket(); err != nil {
		ri.logger.Error("Error stopping WebSocket for %s: %v", ri.exchangeName, err)
	}

	ri.isActive = false
	ri.logger.Info("Stopped %s realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *RealtimeIngestor) IsActive() bool {
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.isActive
}

func (ri *RealtimeIngestor) GetMarketType() connector.MarketType {
	return ri.marketType
}

func (ri *RealtimeIngestor) GetActiveConnections() map[connector.ExchangeName]interface{} {
	ri.mu.RLock()
	defer ri.mu.RUnlock()

	if ri.isActive {
		return map[connector.ExchangeName]interface{}{
			ri.exchangeName: ri.conn,
		}
	}

	return make(map[connector.ExchangeName]interface{})
}

var _ realtime.RealtimeIngestor = (*RealtimeIngestor)(nil)
