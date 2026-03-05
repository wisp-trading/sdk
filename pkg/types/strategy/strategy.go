package strategy

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// StrategyType represents the type of trading strategy
type StrategyType string

const (
	StrategyTypeVolumeMaximizer StrategyType = "volume_maximizer"
	StrategyTypeCashCarry       StrategyType = "cash_carry"
	StrategyTypeArbitrage       StrategyType = "arbitrage"
	StrategyTypeTechnical       StrategyType = "technical"
	StrategyTypeMomentum        StrategyType = "momentum"
	StrategyTypeMeanReversion   StrategyType = "mean_reversion"
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
	GetDescription() string
	GetRiskLevel() RiskLevel
	GetStrategyType() StrategyType

	// Lifecycle — the strategy manages its own execution goroutine and internal clock.
	// Start launches the strategy's run loop. It must be non-blocking.
	Start(ctx context.Context) error
	// Stop signals the strategy to shut down and waits for it to exit cleanly.
	Stop(ctx context.Context) error

	// Signals returns a read-only channel for observing emitted signals.
	// This is an observability tap — production routing goes via wisp.Emit.
	Signals() <-chan Signal
}

// RequiredAsset specifies an asset and which instrument types are needed
type RequiredAsset struct {
	Symbol      portfolio.Pair
	Instruments []connector.Instrument
}

type StrategyConfig struct {
	Name            StrategyName
	Enabled         bool
	MaxPositionSize float64
}

type StrategyMetadata interface {
	GetRequiredFields() []string
	Validate() error
}

// RiskLevel represents the risk classification of a strategy
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

type StrategyExecution struct {
	Orders []connector.Order
	Trades []connector.Trade
}
