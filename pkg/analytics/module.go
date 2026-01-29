package analytics

import (
	"github.com/wisp-trading/wisp/pkg/analytics/analytics"
	"github.com/wisp-trading/wisp/pkg/analytics/indicators"
	"github.com/wisp-trading/wisp/pkg/analytics/market"
	"go.uber.org/fx"
)

var Module = fx.Module("analytics",
	fx.Provide(
		indicators.NewIndicators,
		market.NewMarketService,
		analytics.NewAnalyticsService,
	),
)
