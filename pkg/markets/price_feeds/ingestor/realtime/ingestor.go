package realtime

import (
	"context"
	"sync"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type ingestor struct {
	connector Connector
	store     priceFeedTypes.PriceFeedsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// New creates a new realtime price feed ingestor for WebSocket connectors.
func New(connector Connector, store priceFeedTypes.PriceFeedsStore, logger logging.ApplicationLogger) Ingestor {
	return &ingestor{
		connector: connector,
		store:     store,
		logger:    logger,
	}
}

// Start begins listening to the connector and persisting updates.
func (i *ingestor) Start(ctx context.Context) error {
	if i.isActive {
		return nil
	}

	i.isActive = true
	i.ctx, i.cancel = context.WithCancel(ctx)

	// Subscribe to price updates from the connector
	updateCh := i.connector.Subscribe("")

	// Start processing updates in background
	i.wg.Add(1)
	go i.processUpdates(updateCh)

	return nil
}

// Stop ceases consuming updates and closes the ingestor.
func (i *ingestor) Stop() error {
	if !i.isActive {
		return nil
	}

	i.isActive = false

	if i.cancel != nil {
		i.cancel()
	}

	i.wg.Wait()
	return nil
}

// IsActive returns whether the ingestor is running.
func (i *ingestor) IsActive() bool {
	return i.isActive
}

// GetMarketType returns the market type for price feeds.
func (i *ingestor) GetMarketType() connector.MarketType {
	return connector.MarketType("price_feeds")
}

// GetActiveConnections returns the active connector.
func (i *ingestor) GetActiveConnections() map[connector.ExchangeName]interface{} {
	if !i.isActive {
		return nil
	}
	return map[connector.ExchangeName]interface{}{
		"pyth": i.connector,
	}
}

// processUpdates listens for price updates and persists them to the store.
func (i *ingestor) processUpdates(updateCh <-chan priceFeedTypes.PriceFeedUpdate) {
	defer i.wg.Done()

	for {
		select {
		case <-i.ctx.Done():
			return
		case update := <-updateCh:
			snap := priceFeedTypes.PriceSnapshot{
				FeedID:    update.FeedID,
				Price:     update.Price,
				Timestamp: update.Timestamp,
			}
			if err := i.store.RecordPrice(snap); err != nil {
				i.logger.Errorf("failed to record price for feed %s: %v", update.FeedID, err)
			}
		}
	}
}
