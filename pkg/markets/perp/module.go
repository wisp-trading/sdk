package perp

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/executor"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store"
	"github.com/wisp-trading/sdk/pkg/markets/perp/views"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("perp",
	fx.Provide(
		store.NewStore,
		views.NewPerpViews,
		executor.NewExecutor,
		NewPerpWatchlist,
		NewPerpUniverseProvider,
		fx.Annotate(batch.NewFactory, fx.As(new(batchTypes.BatchIngestorFactory))),
		fx.Annotate(realtime.NewFactory, fx.As(new(realtimeTypes.RealtimeIngestorFactory))),
		fx.Annotate(
			newPerpDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newPerpDomainLifecycle(
	batchFactory batchTypes.BatchIngestorFactory,
	realtimeFactory realtimeTypes.RealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("perp", batchFactory, realtimeFactory, logger)
}
