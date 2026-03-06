package perp

import (
	"github.com/wisp-trading/sdk/pkg/markets/perp/executor"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store"
	"github.com/wisp-trading/sdk/pkg/markets/perp/views"
	"go.uber.org/fx"
)

var Module = fx.Module("perp",
	fx.Provide(
		store.NewStore,
		views.NewPerpViews,
		executor.NewExecutor,
		NewPerpWatchlist,
		NewPerpUniverseProvider,
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
