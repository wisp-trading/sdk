package activity

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PNL provides profit and loss calculations
type PNL interface {
	GetRealizedPNL(ctx context.Context, name strategy.StrategyName) numerical.Decimal
	GetFeesByStrategy(ctx context.Context, name strategy.StrategyName) numerical.Decimal

	// GetRealizedPNLByPair Global realized PNL across all strategies.
	GetRealizedPNLByPair(ctx context.Context, asset portfolio.Pair) numerical.Decimal
	GetTotalRealizedPNL(ctx context.Context) numerical.Decimal

	GetUnrealizedPNL(ctx context.Context, name strategy.StrategyName) (numerical.Decimal, error)
	GetTotalUnrealizedPNL(ctx context.Context) (numerical.Decimal, error)

	// GetTotalPNL Combined global PNL.
	GetTotalPNL(ctx context.Context) (numerical.Decimal, error)

	// GetTotalFees Fee tracking — global.
	GetTotalFees(ctx context.Context) numerical.Decimal
}
