package position

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

func (ds *dataStore) AddOrderToStrategy(strategyName strategy.StrategyName, order connector.Order) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getExecutions()
	execution := current[strategyName]

	if execution == nil {
		execution = &strategy.StrategyExecution{
			Orders:   []connector.Order{},
			TradeIDs: []string{},
		}
	}

	// Add order
	execution.Orders = append(execution.Orders, order)

	// Store updated map
	updated := make(map[strategy.StrategyName]*strategy.StrategyExecution, len(current)+1)
	for k, v := range current {
		updated[k] = v
	}
	updated[strategyName] = execution
	ds.executions.Store(updated)
}

func (ds *dataStore) UpdateOrderInStrategy(strategyName strategy.StrategyName, orderID string, updater func(*connector.Order)) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getExecutions()
	execution := current[strategyName]

	if execution == nil {
		return fmt.Errorf("strategy execution not found: %s", strategyName)
	}

	// Find and update order
	found := false
	for i := range execution.Orders {
		if execution.Orders[i].ID == orderID {
			updater(&execution.Orders[i])
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// Store updated map
	updated := make(map[strategy.StrategyName]*strategy.StrategyExecution, len(current))
	for k, v := range current {
		updated[k] = v
	}
	updated[strategyName] = execution
	ds.executions.Store(updated)

	return nil
}

func (ds *dataStore) CancelOrder(strategyName strategy.StrategyName, orderID string) error {
	return ds.UpdateOrderInStrategy(strategyName, orderID, func(order *connector.Order) {
		order.Status = connector.OrderStatusCancelled
	})
}

func (ds *dataStore) CancelAllPendingOrders(strategyName strategy.StrategyName) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getExecutions()
	execution := current[strategyName]

	if execution == nil {
		return fmt.Errorf("strategy execution not found: %s", strategyName)
	}

	// Cancel all pending orders
	for i := range execution.Orders {
		if execution.Orders[i].Status == connector.OrderStatusPending || execution.Orders[i].Status == connector.OrderStatusOpen {
			execution.Orders[i].Status = connector.OrderStatusCancelled
		}
	}

	// Store updated map
	updated := make(map[strategy.StrategyName]*strategy.StrategyExecution, len(current))
	for k, v := range current {
		updated[k] = v
	}
	updated[strategyName] = execution
	ds.executions.Store(updated)

	return nil
}

func (ds *dataStore) CancelOrdersNotAtLevels(strategyName strategy.StrategyName, validLevels map[string]bool) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getExecutions()
	execution := current[strategyName]

	if execution == nil {
		return fmt.Errorf("strategy execution not found: %s", strategyName)
	}

	// Cancel orders not at valid levels
	for i := range execution.Orders {
		if execution.Orders[i].Status == connector.OrderStatusPending || execution.Orders[i].Status == connector.OrderStatusOpen {
			priceLevel := execution.Orders[i].Price.String()
			if !validLevels[priceLevel] {
				execution.Orders[i].Status = connector.OrderStatusCancelled
			}
		}
	}

	// Store updated map
	updated := make(map[strategy.StrategyName]*strategy.StrategyExecution, len(current))
	for k, v := range current {
		updated[k] = v
	}
	updated[strategyName] = execution
	ds.executions.Store(updated)

	return nil
}
