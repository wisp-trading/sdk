package signal

import (
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// factory is the concrete implementation of wisp.SignalFactory.
type factory struct {
	timeProvider temporal.TimeProvider
}

// NewFactory creates a new signal factory with the injected time provider.
func NewFactory(timeProvider temporal.TimeProvider) strategy.SignalFactory {
	return factory{
		timeProvider: timeProvider,
	}
}

// New creates a new signal builder for a strategy.
func (f factory) New(strategyName strategy.StrategyName) strategy.SignalBuilder {
	return &builder{
		strategyName: strategyName,
		actions:      make([]strategy.TradeAction, 0),
		timeProvider: f.timeProvider,
	}
}
