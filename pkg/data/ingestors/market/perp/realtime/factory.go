package realtime

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/realtime"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

// Factory creates realtime ingestors for all registered perp WebSocket connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.AssetRegistry
	store             perpStore.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.AssetRegistry,
	store perpStore.MarketStore,
	logger logging.ApplicationLogger,
) realtime.RealtimeIngestorFactory {
	return &Factory{
		connectorRegistry: connectorRegistry,
		assetRegistry:     assetRegistry,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered perp WebSocket connector
func (f *Factory) CreateIngestors() []realtime.RealtimeIngestor {
	perpWSConnectors := f.connectorRegistry.GetReadyPerpWebSocketConnectors()

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(perpWSConnectors))

	for _, wsConn := range perpWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		fundingExt := NewFundingRateExtension(f.store, f.logger)

		// Base ingestor + perp extensions
		ingestor := ingestors.NewRealtimeIngestor(
			wsConn, // Perp WebSocket connector
			exchangeName,
			connector.MarketTypePerp,
			f.assetRegistry,
			f.store, // Perp store
			f.logger,
			fundingExt, // Add funding rate WebSocket subscriptions
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created perp realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
