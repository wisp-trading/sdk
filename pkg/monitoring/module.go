package monitoring

import (
	predictionViews "github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	"github.com/wisp-trading/sdk/pkg/monitoring/health"
	"github.com/wisp-trading/sdk/pkg/monitoring/profiling"
	"go.uber.org/fx"
)

var Module = fx.Module("monitoring",
	health.Module,
	profiling.Module,

	fx.Provide(
		NewServer,
		NewViewRegistry,
		predictionViews.NewPredictionViews,
	),
)
