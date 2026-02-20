package real_time

import (
	"context"
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// RealtimeIngestor is a generic base implementation for WebSocket data collection
type RealtimeIngestor struct {
	conn          interface{}
	wsCapable     connector.WebSocketCapable
	exchangeName  connector.ExchangeName
	marketType    connector.MarketType
	assetRegistry registry.PairRegistry
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
	assetRegistry registry.PairRegistry,
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

	return &RealtimeIngestor{
		conn:          conn,
		wsCapable:     wsCapable,
		exchangeName:  exchangeName,
		marketType:    marketType,
		assetRegistry: assetRegistry,
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

	pairs := ri.assetRegistry.GetRequiredPairs()
	if len(pairs) == 0 {
		ri.logger.Warn("No pairs registered for %s realtime ingestion", ri.exchangeName)
		return nil
	}

	ri.ctx, ri.cancel = context.WithCancel(ctx)

	// Start WebSocket connection
	if err := ri.wsCapable.StartWebSocket(); err != nil {
		return fmt.Errorf("failed to start WebSocket for %s: %w", ri.exchangeName, err)
	}

	// All WebSocket subscriptions now handled via extensions
	for _, ext := range ri.extensions {
		if err := ext.Subscribe(ri.conn, ri.exchangeName, pairs); err != nil {
			ri.logger.Error("Failed to subscribe to extension data for %s: %v",
				ri.exchangeName, err)
		}
	}

	// Start extension channel processing
	for _, ext := range ri.extensions {
		go ext.ProcessChannels(ri.conn, ri.exchangeName, ri.ctx)
	}

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

	// Cancel context to stop all extension goroutines
	if ri.cancel != nil {
		ri.cancel()
	}

	// Call unsubscribe on all extensions for cleanup
	for _, ext := range ri.extensions {
		if err := ext.Unsubscribe(ri.conn, ri.exchangeName); err != nil {
			ri.logger.Warn("Failed to unsubscribe extension for %s: %v",
				ri.exchangeName, err)
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
