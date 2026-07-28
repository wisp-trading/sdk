package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// NewOnchainBuilder creates a new onchain signal builder.
func NewOnchainBuilder(strategyName strategy.StrategyName, timeProvider temporal.TimeProvider) *onchainBuilder {
	return &onchainBuilder{
		Core: baseSignal.NewCore(strategyName, timeProvider, connector.MarketTypeOnchain),
	}
}
