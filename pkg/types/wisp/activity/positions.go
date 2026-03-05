package activity

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// Positions provides read-only access to position data
type Positions interface {
	// Strategy-scoped queries
	GetStrategyExecution(name strategy.StrategyName) *strategy.StrategyExecution
	GetTradesForStrategy(name strategy.StrategyName) []connector.Trade

	// Global queries
	GetAllStrategyExecutions() map[strategy.StrategyName]*strategy.StrategyExecution
	GetStrategyForOrder(ctx context.Context, orderID string) (strategy.StrategyName, bool)
	GetTotalOrderCount(ctx context.Context) int64
}
