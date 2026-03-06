package spot

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/spot/executor"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/spot/store"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/markets/spot/views"
	"go.uber.org/fx"
)

// Module wires all spot market dependencies: store, ingestors, views, executor, watchlist, and universe provider.
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

	fx.Invoke(registerStore),
)

func registerStore(
	registry marketTypes.MarketRegistry,
	spotStore domainTypes.MarketStore,
) {
	registry.Register(spotStore)
}
