package perp

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	perpActivity "github.com/wisp-trading/sdk/pkg/markets/perp/activity"
	"github.com/wisp-trading/sdk/pkg/markets/perp/executor"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
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
		perpActivity.NewPerpPNL,
		batch.NewFactory,
		realtime.NewFactory,
		newPerpAssetLoader,
		fx.Annotate(
			newPerpDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newPerpDomainLifecycle(
	assetLoader perpTypes.AssetLoader,
	batchFactory perpTypes.PerpBatchIngestorFactory,
	realtimeFactory perpTypes.PerpRealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("perp", assetLoader, batchFactory, realtimeFactory, logger)
}
