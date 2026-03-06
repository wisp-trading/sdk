package real_time

import (
	"context"
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// realtimeIngestor is a generic base implementation for WebSocket data collection.
type realtimeIngestor struct {
	conn            interface{}
	wsCapable       connector.WebSocketCapable
	exchangeName    connector.ExchangeName
	marketType      connector.MarketType
	marketWatchlist types.MarketWatchlist
	logger          logging.ApplicationLogger

	// State
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
	mu       sync.RWMutex

	// Watchlist subscription
	eventsChan <-chan types.MarketWatchEvent

	// Extension point for market-specific WebSocket subscriptions
	extensions []realtime.WebSocketExtension
}

func NewRealtimeIngestor(
	conn interface{},
	exchangeName connector.ExchangeName,
	marketType connector.MarketType,
	marketWatchlist types.MarketWatchlist,
	logger logging.ApplicationLogger,
	extensions ...realtime.WebSocketExtension,
) realtime.RealtimeIngestor {
	wsCapable, ok := conn.(connector.WebSocketCapable)
	if !ok {
		logger.Error("Connector does not implement WebSocketCapable interface")
		return nil
	}

	return &realtimeIngestor{
		conn:            conn,
		wsCapable:       wsCapable,
		exchangeName:    exchangeName,
		marketType:      marketType,
		marketWatchlist: marketWatchlist,
		logger:          logger,
		extensions:      extensions,
	}
}

func (ri *realtimeIngestor) Start(ctx context.Context) error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if ri.isActive {
		return fmt.Errorf("realtime ingestor for %s already active", ri.exchangeName)
	}

	ri.ctx, ri.cancel = context.WithCancel(ctx)

	// Start WebSocket connection
	if err := ri.wsCapable.StartWebSocket(); err != nil {
		return fmt.Errorf("failed to start WebSocket for %s: %w", ri.exchangeName, err)
	}

	// Initial snapshot of required pairs from the watchlist
	pairs := ri.marketWatchlist.GetRequiredPairs(ri.exchangeName)
	if len(pairs) == 0 {
		ri.logger.Warn("No pairs registered for %s realtime ingestion", ri.exchangeName)
	} else {
		for _, ext := range ri.extensions {
			if err := ext.Subscribe(ri.conn, ri.exchangeName, pairs); err != nil {
				ri.logger.Error("Failed initial subscribe for %s: %v", ri.exchangeName, err)
			}
		}
	}

	// Subscribe to dynamic watchlist events for this exchange
	ch := ri.marketWatchlist.Subscribe(ri.exchangeName)
	ri.eventsChan = ch
	go ri.runWatchlistLoop(ch)

	// Start extension channel processing
	for _, ext := range ri.extensions {
		go ext.ProcessChannels(ri.conn, ri.exchangeName, ri.ctx)
	}

	ri.isActive = true
	ri.logger.Info("Started %s realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *realtimeIngestor) runWatchlistLoop(events <-chan types.MarketWatchEvent) {
	defer ri.marketWatchlist.Unsubscribe(ri.exchangeName)

	for {
		select {
		case <-ri.ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}

			switch ev.Type {
			case types.PairAdded:
				for _, ext := range ri.extensions {
					if err := ext.Subscribe(ri.conn, ri.exchangeName, []portfolio.Pair{ev.Requirement.Pair}); err != nil {
						ri.logger.Error(
							"Failed dynamic subscribe %s %s: %v",
							ri.exchangeName,
							ev.Requirement.Pair.Symbol(),
							err,
						)
					}
				}
			case types.PairRemoved:
				for _, ext := range ri.extensions {
					if err := ext.Unsubscribe(ri.conn, ri.exchangeName); err != nil {
						ri.logger.Warn(
							"Failed dynamic unsubscribe %s %s: %v",
							ri.exchangeName,
							ev.Requirement.Pair.Symbol(),
							err,
						)
					}
				}
			}
		}
	}
}

func (ri *realtimeIngestor) Stop() error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if !ri.isActive {
		return nil
	}

	// Cancel context to stop extension goroutines and watchlist loop
	if ri.cancel != nil {
		ri.cancel()
	}

	if ri.eventsChan != nil {
		ri.marketWatchlist.Unsubscribe(ri.exchangeName)
		ri.eventsChan = nil
	}

	// Stop WebSocket connection
	if err := ri.wsCapable.StopWebSocket(); err != nil {
		ri.logger.Error("Error stopping WebSocket for %s: %v", ri.exchangeName, err)
	}

	ri.isActive = false
	ri.logger.Info("Stopped %s realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *realtimeIngestor) IsActive() bool {
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.isActive
}

func (ri *realtimeIngestor) GetMarketType() connector.MarketType {
	return ri.marketType
}

func (ri *realtimeIngestor) GetActiveConnections() map[connector.ExchangeName]interface{} {
	ri.mu.RLock()
	defer ri.mu.RUnlock()

	if ri.isActive {
		return map[connector.ExchangeName]interface{}{
			ri.exchangeName: ri.conn,
		}
	}

	return make(map[connector.ExchangeName]interface{})
}

var _ realtime.RealtimeIngestor = (*realtimeIngestor)(nil)
