package analytics

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/market"
	"go.uber.org/fx"
)

var Module = fx.Module("analytics",
	fx.Provide(
		indicators.NewIndicators,
		market.NewMarketService,
		analytics.NewAnalyticsService,
	),
)
