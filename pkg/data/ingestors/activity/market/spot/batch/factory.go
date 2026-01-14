package batch

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/batch"
	spotStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/spot"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// Factory creates batch ingestors for all registered spot connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.AssetRegistry
	store             spotStore.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.AssetRegistry,
	store spotStore.MarketStore,
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

// CreateIngestors creates one batch ingestor per registered spot connector
func (f *Factory) CreateIngestors() []batch.BatchIngestor {
	spotConnectors := f.connectorRegistry.GetReadySpotConnectors()

	ingestorList := make([]batch.BatchIngestor, 0, len(spotConnectors))

	for _, conn := range spotConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := ingestors.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypeSpot,
			f.assetRegistry,
			f.store,
			f.timeProvider,
			f.logger,
		)

		ingestorList = append(ingestorList, ingestor)

		f.logger.Info("Created spot batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
