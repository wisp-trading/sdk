package market

import (
	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
)

func (ds *dataStore) GetLastUpdated() marketTypes.LastUpdatedMap {
	return ds.getLastUpdated()
}

func (ds *dataStore) UpdateLastUpdated(key marketTypes.UpdateKey) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getLastUpdated()
	updated := make(marketTypes.LastUpdatedMap, len(current)+1)
	for k, v := range current {
		updated[k] = v
	}
	updated[key] = ds.timeProvider.Now()
	ds.lastUpdated.Store(updated)
}
