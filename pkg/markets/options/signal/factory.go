package signal

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewOptionsBuilder creates a new options signal builder.
func NewOptionsBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) optionsTypes.OptionsSignalBuilder {
	return &optionsBuilder{
		strategyName: strategyName,
		actions:      make([]optionsTypes.OptionsAction, 0),
		timeProvider: timeProvider,
	}
}
