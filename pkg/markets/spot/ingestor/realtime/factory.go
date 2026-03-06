package realtime

import (
	baseRealtime "github.com/wisp-trading/sdk/pkg/markets/base/ingestor/realtime"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered spot WebSocket connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	watchlist         spotTypes.SpotWatchlist
	store             spotTypes.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	watchlist spotTypes.SpotWatchlist,
	store spotTypes.MarketStore,
	logger logging.ApplicationLogger,
) spotTypes.SpotRealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered spot WebSocket connector
func (f *factory) CreateIngestors() []realtimeTypes.RealtimeIngestor {
	spotWSConnectors := f.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())
	realtimeIngestors := make([]realtimeTypes.RealtimeIngestor, 0, len(spotWSConnectors))

	for _, wsConn := range spotWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		// Create extensions for spot markets
		obExt := baseRealtime.NewOrderBookExtension(f.store, f.logger)
		klineExt := baseRealtime.NewKlineExtension(f.store, f.logger, []string{"1m", "5m", "15m", "1h"})

		ingestor := baseRealtime.NewRealtimeIngestor(
			wsConn, exchangeName, connector.MarketTypeSpot,
			f.watchlist, f.logger,
			obExt, klineExt,
		)
		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created spot realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
