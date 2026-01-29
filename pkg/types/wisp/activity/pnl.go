package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PNL provides profit and loss calculations
type PNL interface {
	// Realized PNL (from executed trades)
	GetRealizedPNL(ctx strategy.StrategyContext, strategy strategy.StrategyName) numerical.Decimal
	GetRealizedPNLByAsset(ctx strategy.StrategyContext, asset portfolio.Asset) numerical.Decimal
	GetTotalRealizedPNL(ctx strategy.StrategyContext) numerical.Decimal

	// Unrealized PNL (requires current market prices)
	GetUnrealizedPNL(ctx strategy.StrategyContext, strategy strategy.StrategyName) (numerical.Decimal, error)
	GetTotalUnrealizedPNL(ctx strategy.StrategyContext) (numerical.Decimal, error)

	// Combined
	GetTotalPNL(ctx strategy.StrategyContext) (numerical.Decimal, error)

	// Fee tracking
	GetTotalFees(ctx strategy.StrategyContext) numerical.Decimal
	GetFeesByStrategy(ctx strategy.StrategyContext, strategy strategy.StrategyName) numerical.Decimal
}
