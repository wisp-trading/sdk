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

// NewSpot creates a new spot signal builder for a strategy.
func (f factory) NewSpot(strategyName strategy.StrategyName) strategy.SpotSignalBuilder {
	return &spotBuilder{
		strategyName: strategyName,
		actions:      make([]*strategy.SpotAction, 0),
		timeProvider: f.timeProvider,
	}
}

// NewPerp creates a new perpetual futures signal builder for a strategy.
func (f factory) NewPerp(strategyName strategy.StrategyName) strategy.PerpSignalBuilder {
	return &perpBuilder{
		strategyName: strategyName,
		actions:      make([]*strategy.PerpAction, 0),
		timeProvider: f.timeProvider,
	}
}

// NewPrediction creates a new prediction market signal builder for a strategy.
func (f factory) NewPrediction(strategyName strategy.StrategyName) strategy.PredictionSignalBuilder {
	return &predictionBuilder{
		strategyName: strategyName,
		actions:      make([]*strategy.PredictionAction, 0),
		timeProvider: f.timeProvider,
	}
}

// New creates a new spot signal builder. Kept for backward compatibility.
// Prefer NewSpot for new code.
func (f factory) New(strategyName strategy.StrategyName) strategy.SpotSignalBuilder {
	return f.NewSpot(strategyName)
}
