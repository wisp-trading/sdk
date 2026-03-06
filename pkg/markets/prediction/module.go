package prediction

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/executor"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/signal"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	"go.uber.org/fx"
)

var Module = fx.Module("prediction",
	fx.Provide(
		store.NewStore,
		signal.NewFactory,
		views.NewPredictionViews,
		executor.NewExecutor,
		NewPredictionWatchlist,
		NewPredictionUniverseProvider,
	),

	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ResultTags(`group:"batch_factories"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),
)
