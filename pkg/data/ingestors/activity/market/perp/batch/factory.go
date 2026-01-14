package batch

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/batch"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// Factory creates batch ingestors for all registered perp connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.AssetRegistry
	store             perpStore.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.AssetRegistry,
	store perpStore.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) batch.BatchIngestorFactory {
	return &Factory{
		connectorRegistry: connectorRegistry,
		assetRegistry:     assetRegistry,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

// CreateIngestors creates one batch ingestor per registered perp connector
func (f *Factory) CreateIngestors() []batch.BatchIngestor {
	perpConnectors := f.connectorRegistry.GetReadyPerpConnectors()

	ingestorList := make([]batch.BatchIngestor, 0, len(perpConnectors))

	for _, conn := range perpConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		// Create perp-specific extensions
		fundingExt := NewFundingRateExtension(f.store, f.logger)

		// Base ingestor + perp extensions
		ingestor := ingestors.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypePerp,
			f.assetRegistry,
			f.store,
			f.timeProvider,
			f.logger,
			fundingExt, // Add funding rate collection
		)

		ingestorList = append(ingestorList, ingestor)

		f.logger.Info("Created perp batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
