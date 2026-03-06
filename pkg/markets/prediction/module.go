package prediction

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/executor"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/signal"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/views"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
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
		fx.Annotate(batch.NewFactory, fx.As(new(batchTypes.BatchIngestorFactory))),
		fx.Annotate(realtime.NewFactory, fx.As(new(realtimeTypes.RealtimeIngestorFactory))),
		fx.Annotate(
			newPredictionDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newPredictionDomainLifecycle(
	batchFactory batchTypes.BatchIngestorFactory,
	realtimeFactory realtimeTypes.RealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("prediction", batchFactory, realtimeFactory, logger)
}
