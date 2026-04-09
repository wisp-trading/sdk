package realtime

import (
	"time"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type ingestor struct {
	connector Connector
	store     priceFeedTypes.PriceFeedsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ticker    *time.Ticker
	stopChan  chan struct{}
}

// New creates a new Pyth price feed ingestor following the standard BatchIngestor pattern.
func New(connector Connector, store priceFeedTypes.PriceFeedsStore, logger logging.ApplicationLogger) Ingestor {
	return &ingestor{
		connector: connector,
		store:     store,
		logger:    logger,
		stopChan:  make(chan struct{}),
	}
}

// Start begins polling the connector at the given interval.
func (i *ingestor) Start(interval time.Duration) error {
	if i.isActive {
		return nil
	}

	i.isActive = true
	i.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-i.ticker.C:
				i.CollectNow()
			case <-i.stopChan:
				i.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

// Stop ceases polling and closes the ingestor.
func (i *ingestor) Stop() error {
	if !i.isActive {
		return nil
	}

	i.isActive = false
	close(i.stopChan)
	return nil
}

// IsActive returns whether the ingestor is running.
func (i *ingestor) IsActive() bool {
	return i.isActive
}

// CollectNow fetches and persists price updates.
func (i *ingestor) CollectNow() {
	// Subscribe to price updates from Pyth connector
	updateCh := i.connector.Subscribe("pyth")

	// Collect all available updates (non-blocking)
	for {
		select {
		case update := <-updateCh:
			snap := priceFeedTypes.PriceSnapshot{
				FeedID:    update.FeedID,
				Price:     update.Price,
				Timestamp: update.Timestamp,
			}
			if err := i.store.RecordPrice(snap); err != nil {
				i.logger.Errorf("failed to record price for feed %s: %v", update.FeedID, err)
			}
		default:
			return
		}
	}
}

// GetMarketType returns the market type for price feeds.
func (i *ingestor) GetMarketType() connector.MarketType {
	return connector.MarketType("price_feeds")
}
