package signal

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// factory is the concrete implementation of wisp.SignalFactory.
type factory struct {
	timeProvider temporal.TimeProvider
}

// NewFactory creates a new signal factory with the injected time provider.
func NewFactory(timeProvider temporal.TimeProvider) types.SignalFactory {
	return factory{
		timeProvider: timeProvider,
	}
}

// NewPrediction creates a new prediction market signal builder for a strategy.
func (f factory) NewPrediction(strategyName strategy.StrategyName) types.PredictionSignalBuilder {
	return &predictionBuilder{
		strategyName: strategyName,
		actions:      make([]*types.PredictionAction, 0),
		timeProvider: f.timeProvider,
	}
}
