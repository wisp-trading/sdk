package signal

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewPerpBuilder creates a new perp signal builder.
func NewPerpBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) perpTypes.PerpSignalBuilder {
	return &perpBuilder{
		strategyName: strategyName,
		actions:      make([]perpTypes.PerpAction, 0),
		timeProvider: timeProvider,
	}
}
