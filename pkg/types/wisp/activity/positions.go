package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// Positions provides read-only access to position data
type Positions interface {
	// Strategy queries
	GetStrategyExecution(ctx strategy.StrategyContext) *strategy.StrategyExecution
	GetAllStrategyExecutions(ctx strategy.StrategyContext) map[strategy.StrategyName]*strategy.StrategyExecution

	// Order queries
	GetStrategyForOrder(ctx strategy.StrategyContext, orderID string) (strategy.StrategyName, bool)
	GetTotalOrderCount(ctx strategy.StrategyContext) int64

	// Trade queries
	GetTradesForStrategy(ctx strategy.StrategyContext) []connector.Trade
}
