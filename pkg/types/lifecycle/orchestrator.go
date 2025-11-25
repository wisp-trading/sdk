package lifecycle

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// Orchestrator manages strategy execution lifecycle
type Orchestrator interface {
	// Start begins orchestration
	Start(ctx context.Context) error

	// Stop gracefully stops orchestration
	Stop(ctx context.Context) error

	// NotifyDataUpdate triggers strategy execution on new market data (implements DataUpdateListener)
	NotifyDataUpdate()

	// AddStrategy registers a strategy for execution
	AddStrategy(strat strategy.Strategy)

	// GetStrategies returns all registered strategies
	GetStrategies() []strategy.Strategy
}
