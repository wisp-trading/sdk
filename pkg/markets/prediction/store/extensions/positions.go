package extensions

import (
	"fmt"
	"sync"
	"time"

	predictionTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type predictionPositionsExtension struct {
	mu     sync.RWMutex
	orders map[string]predictionTypes.PredictionOrder // orderID → order
}

// NewPredictionPositionsExtension returns a new PositionsStoreExtension.
func NewPredictionPositionsExtension() predictionTypes.PositionsStoreExtension {
	return &predictionPositionsExtension{
		orders: make(map[string]predictionTypes.PredictionOrder),
	}
}

func (e *predictionPositionsExtension) AddOrder(order predictionTypes.PredictionOrder) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.orders[order.ID] = order
}

// GetOrders returns all orders across all exchanges.
func (e *predictionPositionsExtension) GetOrders() []predictionTypes.PredictionOrder {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]predictionTypes.PredictionOrder, 0, len(e.orders))
	for _, order := range e.orders {
		result = append(result, order)
	}
	return result
}

// GetOrdersByExchange returns all orders placed on a specific exchange.
func (e *predictionPositionsExtension) GetOrdersByExchange(exchange connector.ExchangeName) []predictionTypes.PredictionOrder {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]predictionTypes.PredictionOrder, 0)
	for _, order := range e.orders {
		if order.Exchange == exchange {
			result = append(result, order)
		}
	}
	return result
}

func (e *predictionPositionsExtension) UpdateOrderStatus(orderID string, status connector.OrderStatus) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	order, exists := e.orders[orderID]
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	e.orders[orderID] = order
	return nil
}

func (e *predictionPositionsExtension) UpdateRealizedPnL(orderID string, realized numerical.Decimal) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	order, exists := e.orders[orderID]
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}

	order.RealizedPnL = realized
	order.UpdatedAt = time.Now()
	e.orders[orderID] = order
	return nil
}

func (e *predictionPositionsExtension) QueryOrders(q predictionTypes.PredictionActivityQuery) []predictionTypes.PredictionOrder {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var out []predictionTypes.PredictionOrder
	for _, order := range e.orders {
		if q.Exchange != nil && order.Exchange != *q.Exchange {
			continue
		}
		if q.MarketID != nil && order.MarketID != *q.MarketID {
			continue
		}
		out = append(out, order)
	}
	return out
}

var _ predictionTypes.PositionsStoreExtension = (*predictionPositionsExtension)(nil)
