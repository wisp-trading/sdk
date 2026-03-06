package realtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

// predictionRealtimeIngestor is a WebSocket ingestor for prediction markets.
type predictionRealtimeIngestor struct {
	conn         interface{}
	wsCapable    connector.WebSocketCapable
	exchangeName connector.ExchangeName
	marketType   connector.MarketType
	watchlist    types.PredictionWatchlist
	logger       logging.ApplicationLogger

	// State
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
	mu       sync.RWMutex

	// Watchlist subscription
	eventsChan chan types.PredictionWatchEvent

	// Extension point for prediction-specific WS handling
	extensions []types.PredictionExtension
}

func NewPredictionRealtimeIngestor(
	conn interface{},
	exchangeName connector.ExchangeName,
	marketType connector.MarketType,
	watchlist types.PredictionWatchlist,
	logger logging.ApplicationLogger,
	extensions ...types.PredictionExtension,
) realtime.RealtimeIngestor {
	wsCapable, ok := conn.(connector.WebSocketCapable)
	if !ok {
		logger.Error("Connector does not implement WebSocketCapable interface")
		return nil
	}

	return &predictionRealtimeIngestor{
		conn:         conn,
		wsCapable:    wsCapable,
		exchangeName: exchangeName,
		marketType:   marketType,
		watchlist:    watchlist,
		logger:       logger,
		extensions:   extensions,
	}
}

func (ri *predictionRealtimeIngestor) Start(ctx context.Context) error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if ri.isActive {
		return fmt.Errorf("prediction realtime ingestor for %s already active", ri.exchangeName)
	}

	ri.ctx, ri.cancel = context.WithCancel(ctx)

	if err := ri.wsCapable.StartWebSocket(); err != nil {
		return fmt.Errorf("failed to start WebSocket for %s: %w", ri.exchangeName, err)
	}

	// Initial snapshot: markets to watch for this exchange
	markets := ri.watchlist.GetRequiredMarkets(ri.exchangeName)
	if len(markets) == 0 {
		ri.logger.Warn("No prediction markets registered for %s realtime ingestion", ri.exchangeName)
	} else {
		for _, ext := range ri.extensions {
			for _, market := range markets {
				if err := ext.Subscribe(ri.conn, ri.exchangeName, market); err != nil {
					ri.logger.Error("Failed initial prediction subscribe for %s: %v", ri.exchangeName, err)
				}
			}
		}
	}

	ch := ri.watchlist.Subscribe(ri.exchangeName)
	ri.eventsChan = ch
	go ri.runWatchlistLoop(ch)

	for _, ext := range ri.extensions {
		go ext.ProcessChannels(ri.conn, ri.exchangeName, ri.ctx)
	}

	ri.isActive = true
	ri.logger.Info("Started %s prediction realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *predictionRealtimeIngestor) runWatchlistLoop(events chan types.PredictionWatchEvent) {
	for {
		select {
		case <-ri.ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}

			switch ev.Type {
			case types.PredictionMarketAdded:
				for _, ext := range ri.extensions {
					if err := ext.Subscribe(ri.conn, ri.exchangeName, ev.Market); err != nil {
						ri.logger.Error(
							"Failed dynamic prediction subscribe %s %s: %v",
							ri.exchangeName,
							ev.Market.MarketID,
							err,
						)
					}
				}

			case types.PredictionMarketRemoved:
				for _, ext := range ri.extensions {
					if err := ext.Unsubscribe(ri.conn, ri.exchangeName, ev.Market); err != nil {
						ri.logger.Warn(
							"Failed dynamic prediction unsubscribe %s %s: %v",
							ri.exchangeName,
							ev.Market.MarketID,
							err,
						)
					}
				}
			}
		}
	}
}

func (ri *predictionRealtimeIngestor) Stop() error {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	if !ri.isActive {
		return nil
	}

	if ri.cancel != nil {
		ri.cancel()
	}

	if ri.eventsChan != nil {
		ri.watchlist.Unsubscribe(ri.exchangeName)
		ri.eventsChan = nil
	}

	if err := ri.wsCapable.StopWebSocket(); err != nil {
		ri.logger.Error("Error stopping WebSocket for %s: %v", ri.exchangeName, err)
	}

	ri.isActive = false
	ri.logger.Info("Stopped %s prediction realtime ingestion for %s", ri.marketType, ri.exchangeName)
	return nil
}

func (ri *predictionRealtimeIngestor) IsActive() bool {
	ri.mu.RLock()
	defer ri.mu.RUnlock()
	return ri.isActive
}

func (ri *predictionRealtimeIngestor) GetMarketType() connector.MarketType {
	return ri.marketType
}

func (ri *predictionRealtimeIngestor) GetActiveConnections() map[connector.ExchangeName]interface{} {
	ri.mu.RLock()
	defer ri.mu.RUnlock()

	if ri.isActive {
		return map[connector.ExchangeName]interface{}{
			ri.exchangeName: ri.conn,
		}
	}

	return make(map[connector.ExchangeName]interface{})
}

// Compile-time check
var _ realtime.RealtimeIngestor = (*predictionRealtimeIngestor)(nil)
