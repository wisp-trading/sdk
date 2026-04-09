package batch

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

func (i *ingestor) Stop() error {
	if !i.isActive {
		return nil
	}

	i.isActive = false
	close(i.stopChan)
	return nil
}

func (i *ingestor) IsActive() bool {
	return i.isActive
}

// CollectNow fetches latest prices from the connector and persists them.
func (i *ingestor) CollectNow() {
	updates, err := i.connector.GetLatestPrices()
	if err != nil {
		i.logger.Errorf("failed to fetch latest prices: %v", err)
		return
	}

	for _, update := range updates {
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

func (i *ingestor) GetMarketType() connector.MarketType {
	return connector.MarketType("price_feeds")
}
