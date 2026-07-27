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
// Strategies are self-directed: they own their execution loop and place orders
// only via market-scoped domain Emit:
//
//	wisp.Spot().Emit / Perp().Emit / Predict().Emit / Options().Emit
//
// The orchestrator only manages lifecycle (Start/Stop).
type Strategy interface {
	GetName() StrategyName

	// Start launches the strategy's run loop. It must be non-blocking.
	Start(ctx context.Context) error
	// Stop signals shutdown and waits for a clean exit.
	Stop(ctx context.Context) error

	// LatestStatus returns the most recent status snapshot (for monitoring).
	LatestStatus() StrategyStatus

	// StatusLog returns up to the last 100 status snapshots, oldest-first.
	StatusLog() []StrategyStatus
}

type StrategyExecution struct {
	Orders []connector.Order
	Trades []connector.Trade
}
