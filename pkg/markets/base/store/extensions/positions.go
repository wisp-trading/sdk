package extensions

import (
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type positionsExtension struct {
	mu     sync.RWMutex
	orders []connector.Order
	byID   map[string]int // orderID -> index
}

func NewPositionsExtension() market.PositionsStoreExtension {
	return &positionsExtension{
		orders: []connector.Order{},
		byID:   make(map[string]int),
	}
}

func (e *positionsExtension) AddOrder(order connector.Order) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.byID[order.ID]; exists {
		return
	}
	e.byID[order.ID] = len(e.orders)
	e.orders = append(e.orders, order)
}

func (e *positionsExtension) UpdateOrderStatus(orderID string, status connector.OrderStatus) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	idx, ok := e.byID[orderID]
	if !ok {
		return fmt.Errorf("order %s not found", orderID)
	}
	e.orders[idx].Status = status
	return nil
}

func (e *positionsExtension) GetOrders() []connector.Order {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]connector.Order, len(e.orders))
	copy(out, e.orders)
	return out
}

func (e *positionsExtension) GetTotalOrderCount() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return int64(len(e.orders))
}

func (e *positionsExtension) QueryOrders(q market.ActivityQuery) []connector.Order {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []connector.Order
	for _, o := range e.orders {
		if q.Exchange != nil && o.Exchange != *q.Exchange {
			continue
		}
		if q.Pair != nil && o.Pair.Symbol() != q.Pair.Symbol() {
			continue
		}
		out = append(out, o)
	}
	return out
}

var _ market.PositionsStoreExtension = (*positionsExtension)(nil)
