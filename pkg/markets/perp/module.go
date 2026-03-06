package perp

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/perp/executor"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/markets/perp/views"
	"go.uber.org/fx"
)

// Module wires all perp market dependencies: store, ingestors, views, executor, watchlist, and universe provider.
var Module = fx.Module("perp",
	fx.Provide(
		store.NewStore,
		views.NewPerpViews,
		executor.NewExecutor,
		NewPerpWatchlist,
		NewPerpUniverseProvider,
	),

	// Ingestors
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

	fx.Invoke(
		registerStore,
	),
)

// registerStore registers the perp store with the market registry.
func registerStore(
	registry marketTypes.MarketRegistry,
	perpStore domainTypes.MarketStore,
) {
	registry.Register(perpStore)
}
