package prediction

import (
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/executor"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/signal"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
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

// registerStore registers the prediction store with the market registry.
func registerStore(
	registry marketTypes.MarketRegistry,
	predictionStore domainTypes.MarketStore,
) {
	registry.Register(predictionStore)
}
