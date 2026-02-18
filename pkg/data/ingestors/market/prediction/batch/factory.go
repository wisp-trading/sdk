package batch

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// Factory creates batch ingestors for all registered prediction connectors
type Factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.PairRegistry
	store             prediction.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.PairRegistry,
	store prediction.MarketStore,
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
	predictionConnectors := f.connectorRegistry.FilterPrediction(registry.NewFilter().ReadyOnly().Build())

	ingestorList := make([]batch.BatchIngestor, 0, len(predictionConnectors))

	for _, conn := range predictionConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		// Base ingestor + perp extensions
		ingestor := ingestors.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypePrediction,
			f.assetRegistry,
			f.store,
			f.timeProvider,
			f.logger,
		)

		ingestorList = append(ingestorList, ingestor)

		f.logger.Info("Created perp batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
