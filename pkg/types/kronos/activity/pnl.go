package activity

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// PNL provides profit and loss calculations
type PNL interface {
	// Realized PNL (from executed trades)
	GetRealizedPNL(ctx context.Context, strategy strategy.StrategyName) numerical.Decimal
	GetRealizedPNLByAsset(ctx context.Context, asset portfolio.Asset) numerical.Decimal
	GetTotalRealizedPNL(ctx context.Context) numerical.Decimal

	// Unrealized PNL (requires current market prices)
	GetUnrealizedPNL(ctx context.Context, strategy strategy.StrategyName) (numerical.Decimal, error)
	GetTotalUnrealizedPNL(ctx context.Context) (numerical.Decimal, error)

	// Combined
	GetTotalPNL(ctx context.Context) (numerical.Decimal, error)

	// Fee tracking
	GetTotalFees(ctx context.Context) numerical.Decimal
	GetFeesByStrategy(ctx context.Context, strategy strategy.StrategyName) numerical.Decimal
}
