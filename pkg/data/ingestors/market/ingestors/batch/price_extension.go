package batch

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type PriceExtension struct {
	marketData connector.MarketDataReader
	store      marketTypes.MarketStore
	logger     logging.ApplicationLogger
}

func NewPriceExtension(
	marketData connector.MarketDataReader,
	store marketTypes.MarketStore,
	logger logging.ApplicationLogger,
) batch.CollectionExtension {
	return &PriceExtension{
		marketData: marketData,
		store:      store,
		logger:     logger,
	}
}

// Collect implements batch.CollectionExtension.
func (e *PriceExtension) Collect(conn connector.Connector, exchangeName connector.ExchangeName, assets []portfolio.Pair) {
	if e.marketData == nil {
		return
	}

	var wg sync.WaitGroup

	for _, pair := range assets {
		wg.Add(1)
		go func(p portfolio.Pair) {
			defer wg.Done()

			price, err := e.marketData.FetchPrice(p)
			if err != nil {
				e.logger.Debug("Failed to fetch price for %s on %s: %v", p.Symbol(), exchangeName, err)
				return
			}

			e.store.UpdatePairPrice(p, exchangeName, *price)
			e.store.UpdateLastUpdated(marketTypes.UpdateKey{
				DataType: marketTypes.DataKeyPairPrice,
				Pair:     p,
				Exchange: exchangeName,
			})

			e.logger.Debug(
				"Updated price for %s on %s = %s",
				p.Symbol(),
				exchangeName,
				price.Price.String(),
			)
		}(pair)
	}

	wg.Wait()
}

var _ batch.CollectionExtension = (*PriceExtension)(nil)
