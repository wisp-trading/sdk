package realtime

import (
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
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

// NewFactory creates a new realtime ingestor factory for options markets
func NewFactory(
	connectorRegistry registryTypes.ConnectorRegistry,
	watchlist optionsTypes.OptionsWatchlist,
	store optionsTypes.OptionsStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsRealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

func (f *factory) CreateIngestors() []realtimeTypes.RealtimeIngestor {
	readyConnectors := f.connectorRegistry.FilterOptions(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())
	if len(readyConnectors) == 0 {
		return nil
	}

	ingestors := make([]realtimeTypes.RealtimeIngestor, 0, len(readyConnectors))
	for _, conn := range readyConnectors {
		wsConn, ok := conn.(optionsconnector.WebSocketConnector)
		if !ok {
			continue
		}

		ingestor := &optionsRealtimeIngestor{
			connector: wsConn,
			watchlist: f.watchlist,
			store:     f.store,
			logger:    f.logger,
			isActive:  false,
		}
		ingestors = append(ingestors, ingestor)
	}

	return ingestors
}
