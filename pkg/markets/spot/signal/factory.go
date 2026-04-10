package signal

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewSpotBuilder creates a new spot signal builder.
func NewSpotBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) spotTypes.SpotSignalBuilder {
	return &spotBuilder{
		strategyName: strategyName,
		actions:      make([]spotTypes.SpotAction, 0),
		timeProvider: timeProvider,
	}
}
