package position

import (
	storeActivity "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func NewStore() storeActivity.Positions {
	ds := &dataStore{}
	ds.orders.Store([]connector.Order(nil))
	ds.trades.Store([]connector.Trade(nil))
	ds.ordersByID.Store(make(map[string]int))
	ds.lastUpdated.Store(make(storeActivity.LastUpdatedMap))
	return ds
}

func (ds *dataStore) Clear() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.orders.Store([]connector.Order(nil))
	ds.trades.Store([]connector.Trade(nil))
	ds.ordersByID.Store(make(map[string]int))
	ds.lastUpdated.Store(make(storeActivity.LastUpdatedMap))
}

var _ storeActivity.Positions = (*dataStore)(nil)
