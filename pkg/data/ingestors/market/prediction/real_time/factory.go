package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/real_time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered prediction WebSocket connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	marketWatchlist   data.MarketWatchlist
	store             prediction.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	marketWatchlist data.MarketWatchlist,
	store prediction.MarketStore,
	logger logging.ApplicationLogger,
) realtime.RealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		marketWatchlist:   marketWatchlist,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered prediction WebSocket connector
func (f *factory) CreateIngestors() []realtime.RealtimeIngestor {
	predictionWSConnectors := f.connectorRegistry.FilterPrediction(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(predictionWSConnectors))

	for _, wsConn := range predictionWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := real_time.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypePrediction,
			f.marketWatchlist,
			f.logger,
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created prediction realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
