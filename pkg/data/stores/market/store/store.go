package store

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewStore creates a minimal, market-agnostic base store
func NewStore(timeProvider temporal.TimeProvider, storeExtensions ...marketTypes.StoreExtension) marketTypes.MarketStore {
	ds := &dataStore{
		timeProvider: timeProvider,
		extensions:   storeExtensions,
		prices:       make(map[portfolio.Pair]marketTypes.PriceMap),
		lastUpdated:  make(marketTypes.LastUpdatedMap),
	}

	return ds
}

var _ marketTypes.MarketStore = (*dataStore)(nil)
