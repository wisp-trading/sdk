package analytics

import (
	"github.com/wisp-trading/sdk/pkg/analytics/analytics"
	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
	"github.com/wisp-trading/sdk/pkg/analytics/market"
	"go.uber.org/fx"
)

var Module = fx.Module("analytics",
	fx.Provide(
		indicators.NewIndicators,
		market.NewMarketService,
		analytics.NewAnalyticsService,
	),
)
