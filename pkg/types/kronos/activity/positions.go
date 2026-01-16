package activity

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// Positions provides read-only access to position data
type Positions interface {
	// Strategy queries
	GetStrategyExecution(ctx context.Context) *strategy.StrategyExecution
	GetAllStrategyExecutions(ctx context.Context) map[strategy.StrategyName]*strategy.StrategyExecution

	// Order queries
	GetStrategyForOrder(ctx context.Context, orderID string) (strategy.StrategyName, bool)
	GetTotalOrderCount(ctx context.Context) int64

	// Trade queries
	GetTradesForStrategy(ctx context.Context) []connector.Trade
}
