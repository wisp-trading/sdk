package spot

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	spotActivity "github.com/wisp-trading/sdk/pkg/markets/spot/activity"
	"github.com/wisp-trading/sdk/pkg/markets/spot/executor"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/spot/store"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
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
		spotActivity.NewSpotPositions,
		spotActivity.NewSpotTrades,
		spotActivity.NewSpotPNL,
		batch.NewFactory,
		realtime.NewFactory,
		fx.Annotate(
			newSpotDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newSpotDomainLifecycle(
	batchFactory spotTypes.SpotBatchIngestorFactory,
	realtimeFactory spotTypes.SpotRealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("spot", batchFactory, realtimeFactory, logger)
}
