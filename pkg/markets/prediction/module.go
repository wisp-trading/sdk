package prediction

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/signal"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	predictionTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"go.uber.org/fx"
)

// Module wires all prediction market dependencies: store, ingestors, and views.
var Module = fx.Module("prediction",
	fx.Provide(
		store.NewStore,
		signal.NewFactory,
		views.NewPredictionViews,
		NewPredictionWatchlist,
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

// registerStore registers all market stores with the registry
func registerStore(
	registry marketTypes.MarketRegistry,
	predictionStore predictionTypes.MarketStore,
) {
	registry.Register(predictionStore)
}
