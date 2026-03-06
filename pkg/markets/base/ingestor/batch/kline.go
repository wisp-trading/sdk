package batch

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type klineExtension struct {
	marketData connector.MarketDataReader
	store      marketTypes.KlineStoreExtension
	logger     logging.ApplicationLogger

	intervals   []string
	klineLimits map[string]int
}

func NewKlineExtension(
	marketData connector.MarketDataReader,
	store marketTypes.KlineStoreExtension,
	logger logging.ApplicationLogger,
	intervals []string,
	limits map[string]int,
) batch.CollectionExtension {
	if len(intervals) == 0 {
		intervals = []string{"1m", "5m", "15m", "1h", "4h", "1d"}
	}
	if limits == nil {
		limits = map[string]int{
			"1m":  500,
			"5m":  300,
			"15m": 200,
			"1h":  168,
			"4h":  180,
			"1d":  90,
		}
	}

	return &klineExtension{
		marketData:  marketData,
		store:       store,
		logger:      logger,
		intervals:   intervals,
		klineLimits: limits,
	}
}

// Collect implements batch.CollectionExtension.
func (e *klineExtension) Collect(conn connector.Connector, exchangeName connector.ExchangeName, assets []portfolio.Pair) {
	if e.marketData == nil {
		return
	}

	var wg sync.WaitGroup

	for _, pair := range assets {
		for _, interval := range e.intervals {
			wg.Add(1)
			go func(p portfolio.Pair, iv string) {
				defer wg.Done()

				limit := e.klineLimits[iv]
				if limit == 0 {
					limit = 100
				}

				klines, err := e.marketData.FetchKlines(p, iv, limit)
				if err != nil {
					e.logger.Debug("Failed to fetch %s klines for %s on %s: %v", iv, p.Symbol(), exchangeName, err)
					return
				}

				if len(klines) == 0 {
					return
				}

				for _, kline := range klines {
					e.store.UpdateKline(p, exchangeName, kline)
				}

				e.logger.Debug("Updated %d %s klines for %s on %s", len(klines), iv, p.Symbol(), exchangeName)
			}(pair, interval)
		}
	}

	wg.Wait()
}

var _ batch.CollectionExtension = (*klineExtension)(nil)
