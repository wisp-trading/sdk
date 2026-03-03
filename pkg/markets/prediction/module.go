package prediction

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	"go.uber.org/fx"
)

// Module wires all prediction market dependencies: store, ingestors, and views.
var Module = fx.Module("prediction",
	// Store
	fx.Provide(
		store.NewStore,
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

	// Views (monitoring)
	fx.Provide(
		views.NewPredictionViews,
	),
)

// BatchModule exposes the batch factory separately for when batch ingestion is enabled.
var BatchModule = fx.Provide(
	fx.Annotate(
		batch.NewFactory,
		fx.ResultTags(`group:"batch_factories"`),
	),
)
