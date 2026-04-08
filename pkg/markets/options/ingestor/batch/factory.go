package batch

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type factory struct {
	connectorRegistry registryTypes.ConnectorRegistry
	watchlist         optionsTypes.OptionsWatchlist
	store             optionsTypes.OptionsStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

// NewFactory creates a new batch ingestor factory for options markets
func NewFactory(
	connectorRegistry registryTypes.ConnectorRegistry,
	watchlist optionsTypes.OptionsWatchlist,
	store optionsTypes.OptionsStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) batchTypes.BatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

func (f *factory) CreateIngestors() []batchTypes.BatchIngestor {
	readyConnectors := f.connectorRegistry.FilterOptions(registry.NewFilter().ReadyOnly().Build())
	if len(readyConnectors) == 0 {
		return nil
	}

	ingestors := make([]batchTypes.BatchIngestor, 0, len(readyConnectors))
	for _, conn := range readyConnectors {
		ingestor := &optionsIngestor{
			connector: conn,
			watchlist: f.watchlist,
			store:     f.store,
			logger:    f.logger,
			isActive:  false,
			stopChan:  make(chan struct{}),
		}
		ingestors = append(ingestors, ingestor)
	}

	return ingestors
}
