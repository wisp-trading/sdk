package activity

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
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

	// Order management
	AddOrderToStrategy(strategy strategy.StrategyName, order connector.Order)
	UpdateOrderInStrategy(strategy strategy.StrategyName, orderID string, updater func(*connector.Order)) error

	// Trade linking
	LinkTradeToStrategy(strategy strategy.StrategyName, tradeID string)
	GetTradeIDsForStrategy(strategy strategy.StrategyName) []string

	// Order cancellation
	CancelOrder(strategy strategy.StrategyName, orderID string) error
	CancelAllPendingOrders(strategy strategy.StrategyName) error
	CancelOrdersNotAtLevels(strategy strategy.StrategyName, validLevels map[string]bool) error

	// Position reconciliation (validate computed positions vs exchange positions)
	ReconcilePosition(strategyName strategy.StrategyName, exchangePos connector.Position) error

	// Last updated tracking
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)

	// Clear all data for simulation restart
	Clear()
}
