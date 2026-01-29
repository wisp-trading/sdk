package store

import (
	marketTypes "github.com/wisp-trading/wisp/pkg/types/data/stores/market"
	"github.com/wisp-trading/wisp/pkg/types/temporal"
)

func NewStore(timeProvider temporal.TimeProvider, extensions ...marketTypes.StoreExtension) marketTypes.MarketStore {
	ds := &dataStore{
		timeProvider: timeProvider,
		extensions:   extensions,
	}
	ds.orderBooks.Store(make(assetOrderBooks))
	ds.prices.Store(make(assetPrices))
	ds.klines.Store(make(assetKlines))
	ds.lastUpdated.Store(make(marketTypes.LastUpdatedMap))

	return ds
}

// Helper methods to get typed data from atomic.Value

func (ds *dataStore) getOrderBooks() assetOrderBooks {
	if v := ds.orderBooks.Load(); v != nil {
		return v.(assetOrderBooks)
	}
	return make(assetOrderBooks)
}

func (ds *dataStore) getPrices() assetPrices {
	if v := ds.prices.Load(); v != nil {
		return v.(assetPrices)
	}
	return make(assetPrices)
}

func (ds *dataStore) getKlines() assetKlines {
	if v := ds.klines.Load(); v != nil {
		return v.(assetKlines)
	}
	return make(assetKlines)
}

func (ds *dataStore) getLastUpdated() marketTypes.LastUpdatedMap {
	if v := ds.lastUpdated.Load(); v != nil {
		return v.(marketTypes.LastUpdatedMap)
	}
	return make(marketTypes.LastUpdatedMap)
}

var _ marketTypes.MarketStore = (*dataStore)(nil)
