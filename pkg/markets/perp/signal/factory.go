package signal

import (
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewPerpBuilder creates a new perp signal builder. Consumed by pkg/signal/factory.go
// so that strategy.SignalFactory.NewPerp is wired through the perp domain.
func NewPerpBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) strategy.PerpSignalBuilder {
	return &perpBuilder{
		strategyName: strategyName,
		actions:      make([]*strategy.PerpAction, 0),
		timeProvider: timeProvider,
	}
}
