package store

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
)

func (ds *dataStore) GetLastUpdated() marketTypes.LastUpdatedMap {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return shallow copy to prevent external mutation
	result := make(marketTypes.LastUpdatedMap, len(ds.lastUpdated))
	for k, v := range ds.lastUpdated {
		result[k] = v
	}
	return result
}

func (ds *dataStore) UpdateLastUpdated(key marketTypes.UpdateKey) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.lastUpdated[key] = ds.timeProvider.Now()
}
