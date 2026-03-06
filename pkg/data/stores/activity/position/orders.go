package position

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func (ds *dataStore) AddOrder(order connector.Order) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	orders := append(ds.getOrders(), order)
	index := make(map[string]int, len(orders))
	for i, o := range orders {
		index[o.ID] = i
	}
	ds.orders.Store(orders)
	ds.ordersByID.Store(index)
}

func (ds *dataStore) UpdateOrderStatus(orderID string, status connector.OrderStatus) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	orders := ds.getOrders()
	idx, ok := ds.getOrderIndex()[orderID]
	if !ok {
		return fmt.Errorf("order not found: %s", orderID)
	}

	updated := make([]connector.Order, len(orders))
	copy(updated, orders)
	updated[idx].Status = status
	ds.orders.Store(updated)
	return nil
}

// GetTotalOrderCount returns the total number of orders
func (ds *dataStore) GetTotalOrderCount() int64 {
	return int64(len(ds.getOrders()))
}
