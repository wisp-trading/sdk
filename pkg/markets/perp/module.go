package perp

import (
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/markets/perp/views"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"go.uber.org/fx"
)

// Module wires all perp market dependencies: store, ingestors, views, watchlist, and universe provider.
var Module = fx.Module("perp",
	fx.Provide(
		store.NewStore,
		views.NewPerpViews,
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
