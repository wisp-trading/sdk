package prediction

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/executor"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/signal"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"go.uber.org/fx"
)

// Module wires all prediction market dependencies: store, ingestors, views, and executor.
var Module = fx.Module("prediction",
	fx.Provide(
		store.NewStore,
		signal.NewFactory,
		views.NewPredictionViews,
		executor.NewExecutor,
		NewPredictionWatchlist,
		NewPredictionUniverseProvider,
	),

	// Ingestors
	fx.Provide(
		//fx.Annotate(
		//	batch.NewFactory,
		//	fx.ResultTags(`group:"batch_factories"`),
		//),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),

	fx.Invoke(
		registerStore,
	),
)

// BatchModule exposes the batch factory separately for when batch ingestion is enabled.
var BatchModule = fx.Provide(
	fx.Annotate(
		batch.NewFactory,
		fx.ResultTags(`group:"batch_factories"`),
	),
)

// registerStore registers the prediction store with the market registry.
// Accepts domainTypes.Store (the richer domain interface) which satisfies marketTypes.MarketStore.
func registerStore(
	registry marketTypes.MarketRegistry,
	predictionStore domainTypes.Store,
) {
	registry.Register(predictionStore)
}
