package batch

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	batchTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// factory creates batch ingestors for all registered prediction connectors.
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	marketWatchlist   data.MarketWatchlist
	store             prediction.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	marketWatchlist data.MarketWatchlist,
	store prediction.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) batchTypes.BatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		marketWatchlist:   marketWatchlist,
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

		ingestor := batch.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypePrediction,
			f.marketWatchlist,
			f.store,
			f.timeProvider,
			f.logger,
		)

		ingestorList = append(ingestorList, ingestor)
		f.logger.Info("Created prediction batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
