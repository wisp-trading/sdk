package extensions

import (
	"fmt"
	"sync"
	"time"

	predictionTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type predictionPositionsExtension struct {
	mu     sync.RWMutex
	orders map[string]predictionTypes.PredictionOrder // orderID → order
	index  map[strategy.StrategyName][]string         // strategy → []orderID
}

// NewPredictionPositionsExtension returns a new PositionsStoreExtension.
func NewPredictionPositionsExtension() predictionTypes.PositionsStoreExtension {
	return &predictionPositionsExtension{
		orders: make(map[string]predictionTypes.PredictionOrder),
		index:  make(map[strategy.StrategyName][]string),
	}
}

func (e *predictionPositionsExtension) AddOrder(strat strategy.StrategyName, order predictionTypes.PredictionOrder) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.orders[order.ID] = order
	e.index[strat] = append(e.index[strat], order.ID)
}

func (e *predictionPositionsExtension) GetOrdersByStrategy(strat strategy.StrategyName) []predictionTypes.PredictionOrder {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ids, ok := e.index[strat]
	if !ok {
		return nil
	}

	result := make([]predictionTypes.PredictionOrder, 0, len(ids))
	for _, id := range ids {
		if order, exists := e.orders[id]; exists {
			result = append(result, order)
		}
	}
	return result
}

func (e *predictionPositionsExtension) GetStrategyForOrder(orderID string) (strategy.StrategyName, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	order, exists := e.orders[orderID]
	if !exists {
		return "", false
	}
	return order.StrategyName, true
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
