package spot

import (
	"github.com/wisp-trading/sdk/pkg/markets/spot/executor"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/spot/store"
	"github.com/wisp-trading/sdk/pkg/markets/spot/views"
	"go.uber.org/fx"
)

var Module = fx.Module("spot",
	fx.Provide(
		store.NewStore,
		views.NewSpotViews,
		executor.NewExecutor,
		NewSpotWatchlist,
		NewSpotUniverseProvider,
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
