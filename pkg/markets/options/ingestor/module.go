package ingestor

import (
	"github.com/wisp-trading/sdk/pkg/markets/options/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/options/ingestor/realtime"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"go.uber.org/fx"
)

var Module = fx.Module("options.ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ResultTags(`group:"batch_ingestors"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`group:"realtime_ingestors"`),
		),
	),
)

// ProvideBatchIngestorFactory provides the batch ingestor factory for options
func ProvideBatchIngestorFactory(
	connectorRegistry registryTypes.ConnectorRegistry,
	watchlist optionsTypes.OptionsWatchlist,
	store optionsTypes.OptionsStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsBatchIngestorFactory {
	return batch.NewFactory(connectorRegistry, watchlist, store, timeProvider, logger)
}

// ProvideRealtimeIngestorFactory provides the realtime ingestor factory for options
func ProvideRealtimeIngestorFactory(
	connectorRegistry registryTypes.ConnectorRegistry,
	watchlist optionsTypes.OptionsWatchlist,
	store optionsTypes.OptionsStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsRealtimeIngestorFactory {
	return realtime.NewFactory(connectorRegistry, watchlist, store, timeProvider, logger)
}
