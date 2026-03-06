package activity

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type UpdateKey string

type LastUpdatedMap map[UpdateKey]time.Time
type StrategyExecutionMap map[strategy.StrategyName]*strategy.StrategyExecution

type Positions interface {
	// Strategy execution management
	StoreStrategyExecution(strategy strategy.StrategyName, execution *strategy.StrategyExecution)
	GetStrategyExecution(strategy strategy.StrategyName) *strategy.StrategyExecution
	UpdateStrategyExecution(strategy strategy.StrategyName, updateFunc func(*strategy.StrategyExecution)) error

	// Portfolio queries
	GetAllStrategyExecutions() map[strategy.StrategyName]*strategy.StrategyExecution
	GetTotalOrderCount() int64

	// Order storage
	AddOrderToStrategy(strategy strategy.StrategyName, order connector.Order)
	UpdateOrderStatus(strategy strategy.StrategyName, orderID string, status connector.OrderStatus) error
	GetStrategyForOrder(orderID string) (strategy.StrategyName, bool)

	// Trade storage
	AddTradeToStrategy(strategy strategy.StrategyName, trade connector.Trade)
	GetTradesForStrategy(strategy strategy.StrategyName) []connector.Trade

	// Last updated tracking
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)

	// Clear all data for simulation restart
	Clear()
}
