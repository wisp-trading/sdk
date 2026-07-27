package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewPerpBuilder creates a new perp signal builder (shared pair core + leverage extensions).
func NewPerpBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) perpTypes.PerpSignalBuilder {
	return &perpBuilder{
		Core: baseSignal.NewCore(strategyName, timeProvider, connector.MarketTypePerp),
	}
}
