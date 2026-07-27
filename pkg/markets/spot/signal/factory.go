package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewSpotBuilder creates a new spot signal builder (shared pair core + spot MarketType).
func NewSpotBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) spotTypes.SpotSignalBuilder {
	return &spotBuilder{
		Core: baseSignal.NewCore(strategyName, timeProvider, connector.MarketTypeSpot),
	}
}
