package strategy

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type StrategyName string

const (
	CashCarry       StrategyName = "Cash Carry"
	VolumeMaximizer StrategyName = "Volume Maximizer"
	Momentum        StrategyName = "Momentum"
)

// Strategy is the interface that all trading strategies must implement.
// Strategies are self-directed: they own their execution loop and push signals
// asynchronously via wisp.Emit(signal). The orchestrator only manages lifecycle.
type Strategy interface {
	// Identity
	GetName() StrategyName

	// Lifecycle — the strategy manages its own execution goroutine and internal clock.
	// Start launches the strategy's run loop. It must be non-blocking.
	Start(ctx context.Context) error
	// Stop signals the strategy to shut down and waits for it to exit cleanly.
	Stop(ctx context.Context) error

	// Signals returns a read-only channel for observing emitted signals.
	// This is an observability tap — production routing goes via wisp.Emit.
	Signals() <-chan Signal

	// LatestStatus returns the most recently emitted status snapshot.
	// Returns a zero value if no status has been emitted yet.
	LatestStatus() StrategyStatus

	// StatusLog returns up to the last 100 status snapshots, oldest-first.
	// The strategy owns this log; any caller can read it at any time.
	StatusLog() []StrategyStatus
}

type StrategyExecution struct {
	Orders []connector.Order
	Trades []connector.Trade
}
