package store

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewStore creates a minimal, market-agnostic base store
func NewStore(timeProvider temporal.TimeProvider, storeExtensions ...marketTypes.StoreExtension) marketTypes.MarketStore {
	ds := &dataStore{
		timeProvider: timeProvider,
		extensions:   storeExtensions,
	}

	ds.prices.Store(make(assetPrices))
	ds.lastUpdated.Store(make(marketTypes.LastUpdatedMap))

	return ds
}

// Helper methods to get typed data from atomic.Value

func (ds *dataStore) getPrices() assetPrices {
	if v := ds.prices.Load(); v != nil {
		return v.(assetPrices)
	}
	return make(assetPrices)
}

func (ds *dataStore) getLastUpdated() marketTypes.LastUpdatedMap {
	if v := ds.lastUpdated.Load(); v != nil {
		return v.(marketTypes.LastUpdatedMap)
	}
	return make(marketTypes.LastUpdatedMap)
}

var _ marketTypes.MarketStore = (*dataStore)(nil)
