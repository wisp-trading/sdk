package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	spotStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// Factory creates realtime ingestors for all registered spot WebSocket connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.PairRegistry
	store             spotStore.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.PairRegistry,
	store spotStore.MarketStore,
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
	spotWSConnectors := f.connectorRegistry.GetReadySpotWebSocketConnectors()

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(spotWSConnectors))

	for _, wsConn := range spotWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := ingestors.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypeSpot,
			f.assetRegistry,
			f.store,
			f.logger,
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created spot realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
