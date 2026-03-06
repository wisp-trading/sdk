package spot

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/spot/executor"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/spot/store"
	"github.com/wisp-trading/sdk/pkg/markets/spot/views"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("spot",
	fx.Provide(
		store.NewStore,
		views.NewSpotViews,
		executor.NewExecutor,
		NewSpotWatchlist,
		NewSpotUniverseProvider,
		// Ingestor factories — consumed by the domain coordinator below.
		fx.Annotate(batch.NewFactory, fx.As(new(batchTypes.BatchIngestorFactory))),
		fx.Annotate(realtime.NewFactory, fx.As(new(realtimeTypes.RealtimeIngestorFactory))),
		fx.Annotate(
			newSpotDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newSpotDomainLifecycle(
	batchFactory batchTypes.BatchIngestorFactory,
	realtimeFactory realtimeTypes.RealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("spot", batchFactory, realtimeFactory, logger)
}
