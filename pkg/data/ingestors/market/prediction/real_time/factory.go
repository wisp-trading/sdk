package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// Factory creates realtime ingestors for all registered prediction WebSocket connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.PairRegistry
	store             prediction.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.PairRegistry,
	store prediction.MarketStore,
	logger logging.ApplicationLogger,
) realtime.RealtimeIngestorFactory {
	return &Factory{
		connectorRegistry: connectorRegistry,
		assetRegistry:     assetRegistry,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered spot WebSocket connector
func (f *Factory) CreateIngestors() []realtime.RealtimeIngestor {
	predictionWSConnectors := f.connectorRegistry.FilterPrediction(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(predictionWSConnectors))

	for _, wsConn := range predictionWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := ingestors.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypePrediction,
			f.assetRegistry,
			f.store,
			f.logger,
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created spot realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
