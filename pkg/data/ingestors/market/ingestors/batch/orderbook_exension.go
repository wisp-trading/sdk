package batch

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	batchTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type orderBookExtension struct {
	marketData connector.MarketDataReader
	store      marketTypes.OrderBookStoreExtension
	logger     logging.ApplicationLogger
	depth      int
}

func NewOrderBookExtension(
	marketData connector.MarketDataReader,
	store marketTypes.OrderBookStoreExtension,
	logger logging.ApplicationLogger,
	depth int,
) batchTypes.CollectionExtension {
	if depth <= 0 {
		depth = 20
	}

	return &orderBookExtension{
		marketData: marketData,
		store:      store,
		logger:     logger,
		depth:      depth,
	}
}

func (e *orderBookExtension) Collect(
	conn connector.Connector,
	exchangeName connector.ExchangeName,
	assets []portfolio.Pair,
) {
	if e.marketData == nil {
		// Try to recover from conn if not provided
		if md, ok := conn.(connector.MarketDataReader); ok {
			e.marketData = md
		} else {
			return
		}
	}

	var wg sync.WaitGroup

	for _, asset := range assets {
		wg.Add(1)
		go func(a portfolio.Pair) {
			defer wg.Done()

			orderBook, err := e.marketData.FetchOrderBook(a, e.depth)
			if err != nil {
				e.logger.Debug("Failed to fetch order book for %s on %s: %v", a.Symbol(), exchangeName, err)
				return
			}

			if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
				e.logger.Debug("Empty order book for %s on %s", a.Symbol(), exchangeName)
				return
			}

			e.store.UpdateOrderBook(a, exchangeName, *orderBook)

			e.logger.Debug("Updated order book for %s on %s - bid: %s, ask: %s",
				a.Symbol(), exchangeName,
				orderBook.Bids[0].Price.StringFixed(2),
				orderBook.Asks[0].Price.StringFixed(2),
			)
		}(asset)
	}

	wg.Wait()
}

var _ batchTypes.CollectionExtension = (*orderBookExtension)(nil)
