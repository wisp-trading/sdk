package realtime

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type factory struct {
	connectorRegistry registry.ConnectorRegistry
	store             priceFeedTypes.PriceFeedsStore
	logger            logging.ApplicationLogger
}

// NewFactory creates a new batch ingestor factory for price feeds.
func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	store priceFeedTypes.PriceFeedsStore,
	logger logging.ApplicationLogger,
) priceFeedTypes.PriceFeedsBatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		store:             store,
		logger:            logger,
	}
}

func (f *factory) CreateIngestors() []batchTypes.BatchIngestor {
	readyConnectors := f.connectorRegistry.Filter(registry.NewFilter().ReadyOnly().Build())
	if len(readyConnectors) == 0 {
		return nil
	}

	ingestors := make([]batchTypes.BatchIngestor, 0)
	for _, conn := range readyConnectors {
		// Only process Pyth connectors (price feed connectors)
		if pythConn, ok := conn.(Connector); ok {
			ingestor := &ingestor{
				connector: pythConn,
				store:     f.store,
				logger:    f.logger,
				isActive:  false,
				stopChan:  make(chan struct{}),
			}
			ingestors = append(ingestors, ingestor)
		}
	}

	return ingestors
}
