package realtime

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered prediction WebSocket connectors.
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	watchlist         types.PredictionWatchlist
	logger            logging.ApplicationLogger
	store             types.MarketStore
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	watchlist types.PredictionWatchlist,
	logger logging.ApplicationLogger,
	store types.MarketStore,
) realtime.RealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		logger:            logger,
		store:             store,
	}
}

// CreateIngestors creates one realtime ingestor per registered prediction WebSocket connector.
func (f *factory) CreateIngestors() []realtime.RealtimeIngestor {
	predictionWSConnectors := f.connectorRegistry.FilterPrediction(
		registry.NewFilter().ReadyOnly().WebSocketOnly().Build(),
	)

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(predictionWSConnectors))

	for _, wsConn := range predictionWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := NewPredictionRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypePrediction,
			f.watchlist,
			f.logger,
			NewPredictionOrderBookExtension(f.store, f.logger),
		)

		if ingestor == nil {
			continue
		}

		realtimeIngestors = append(realtimeIngestors, ingestor)
		f.logger.Info("Created prediction realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
