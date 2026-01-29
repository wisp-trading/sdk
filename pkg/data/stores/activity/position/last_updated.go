package position

import (
	portfolioTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/activity"
)

func (ds *dataStore) GetLastUpdated() portfolioTypes.LastUpdatedMap {
	if v := ds.lastUpdated.Load(); v != nil {
		return v.(portfolioTypes.LastUpdatedMap)
	}
	return make(portfolioTypes.LastUpdatedMap)
}

func (ds *dataStore) UpdateLastUpdated(key portfolioTypes.UpdateKey) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.GetLastUpdated()
	updated := make(portfolioTypes.LastUpdatedMap, len(current)+1)
	for k, v := range current {
		updated[k] = v
	}
	updated[key] = ds.timeProvider.Now()
	ds.lastUpdated.Store(updated)
}
