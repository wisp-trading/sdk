package position

import (
	"sync"
	"sync/atomic"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type dataStore struct {
	mutex       sync.RWMutex
	orders      atomic.Value // []connector.Order
	trades      atomic.Value // []connector.Trade
	ordersByID  atomic.Value // map[string]int
	lastUpdated atomic.Value // storeActivity.LastUpdatedMap
}

func (ds *dataStore) getOrders() []connector.Order {
	if v := ds.orders.Load(); v != nil {
		return v.([]connector.Order)
	}
	return nil
}

func (ds *dataStore) getTrades() []connector.Trade {
	if v := ds.trades.Load(); v != nil {
		return v.([]connector.Trade)
	}
	return nil
}

func (ds *dataStore) getOrderIndex() map[string]int {
	if v := ds.ordersByID.Load(); v != nil {
		return v.(map[string]int)
	}
	return make(map[string]int)
}
