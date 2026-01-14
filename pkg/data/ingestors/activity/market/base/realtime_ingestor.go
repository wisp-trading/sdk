package base

import (
	"context"
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/common"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

// WebSocketExtension allows market-specific WebSocket subscriptions (funding rate updates, etc.)
type WebSocketExtension interface {
	Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Asset) error
	Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error
}

// RealtimeIngestor is a generic base implementation for WebSocket data collection
// Works with ANY connector that implements WebSocket capabilities
type RealtimeIngestor struct {
	conn          common.BaseConnector
	wsCapable     common.WebSocketCapable
	exchangeName  connector.ExchangeName
	marketType    connector.MarketType
	assetRegistry registry.AssetRegistry
	store         interface{}
	logger        logging.ApplicationLogger

	// State
	isActive bool
	mu       sync.RWMutex

	// Extension point for market-specific WebSocket subscriptions
	extensions []WebSocketExtension
}

func NewRealtimeIngestor(
	conn common.BaseConnector,
	exchangeName connector.ExchangeName,
	marketType connector.MarketType,
	assetRegistry registry.AssetRegistry,
	store interface{},
	logger logging.ApplicationLogger,
	extensions ...WebSocketExtension,
) ingestors.RealtimeIngestor {
	// Cast to WebSocketCapable for lifecycle management
	wsCapable, ok := conn.(common.WebSocketCapable)
	if !ok {
		logger.Error("Connector does not implement WebSocketCapable interface")
		return nil
	}

	return &RealtimeIngestor{
		conn:          conn,
		wsCapable:     wsCapable,
		exchangeName:  exchangeName,
		marketType:    marketType,
		assetRegistry: assetRegistry,
		store:         store,
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

	// Start WebSocket connection
	if err := ri.wsCapable.StartWebSocket(); err != nil {
		return fmt.Errorf("failed to start WebSocket for %s: %w", ri.exchangeName, err)
	}

	// Subscribe to market data via extensions
	// Each market type (spot/perp) handles its own subscriptions
	for _, ext := range ri.extensions {
		if err := ext.Subscribe(ri.conn, ri.exchangeName, assets); err != nil {
			ri.logger.Error("Failed to subscribe to %s-specific data for %s: %v", ri.marketType, ri.exchangeName, err)
			// Continue - not fatal
		}
	}

	// Note: WebSocket data is processed by the connector implementation
	// The connector writes directly to stores or provides channels that we listen to
	// This is handled by market-specific extensions

	ri.isActive = true
	ri.logger.Info("Started %s realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *RealtimeIngestor) Stop() error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if !ri.isActive {
		return nil
	}

	// Unsubscribe from market-specific extensions
	for _, ext := range ri.extensions {
		if err := ext.Unsubscribe(ri.conn, ri.exchangeName); err != nil {
			ri.logger.Warn("Failed to unsubscribe from %s-specific data for %s: %v", ri.marketType, ri.exchangeName, err)
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

var _ ingestors.RealtimeIngestor = (*RealtimeIngestor)(nil)
