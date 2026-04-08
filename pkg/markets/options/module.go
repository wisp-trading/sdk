package options

import (
	baseIngestor "github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	"github.com/wisp-trading/sdk/pkg/markets/options/activity"
	"github.com/wisp-trading/sdk/pkg/markets/options/analytics"
	"github.com/wisp-trading/sdk/pkg/markets/options/executor"
	"github.com/wisp-trading/sdk/pkg/markets/options/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/options/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/options/store"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/markets/options/views"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("options",
	fx.Provide(
		store.NewStore,
		views.NewView,
		executor.NewExecutor,
		NewOptionsWatchlist,
		NewOptionsUniverseProvider,
		activity.NewPNLCalculator,
		analytics.NewAnalyticsService,
		batch.NewFactory,
		realtime.NewFactory,
		NewAssetLoader,
		fx.Annotate(
			newOptionsDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newOptionsDomainLifecycle(
	assetLoader optionsTypes.OptionsAssetLoader,
	batchFactory optionsTypes.OptionsBatchIngestorFactory,
	realtimeFactory optionsTypes.OptionsRealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return baseIngestor.NewDomainCoordinator("options", assetLoader, batchFactory, realtimeFactory, logger)
}
