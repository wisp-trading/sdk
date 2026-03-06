package batch

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// factory creates batch ingestors for all registered prediction connectors.
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	store             types.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	store types.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) batchTypes.BatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

// CreateIngestors creates one batch ingestor per registered prediction connector.
func (f *factory) CreateIngestors() []batchTypes.BatchIngestor {
	predictionConnectors := f.connectorRegistry.FilterPrediction(registry.NewFilter().ReadyOnly().Build())

	ingestorList := make([]batchTypes.BatchIngestor, 0, len(predictionConnectors))

	for _, conn := range predictionConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		ingestor := NewPredictionBatchIngestor(
			conn,
			exchangeName,
			f.logger,
			f.timeProvider,
			NewBalanceCollectionExtension(f.store, f.logger),
		)

		ingestorList = append(ingestorList, ingestor)
		f.logger.Info("Created prediction batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
