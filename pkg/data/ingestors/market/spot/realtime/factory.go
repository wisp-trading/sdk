package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/real_time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	spotStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered spot WebSocket connectors
type factory struct {
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
	return &factory{
		connectorRegistry: connectorRegistry,
		assetRegistry:     assetRegistry,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered spot WebSocket connector
func (f *factory) CreateIngestors() []realtime.RealtimeIngestor {
	spotWSConnectors := f.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(spotWSConnectors))

	for _, wsConn := range spotWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		// Create extensions for spot markets
		obExt := real_time.NewOrderBookExtension(f.store, f.logger)
		klineExt := real_time.NewKlineExtension(f.store, f.logger, []string{"1m", "5m", "15m", "1h"})

		ingestor := real_time.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypeSpot,
			f.assetRegistry,
			f.store,
			f.logger,
			obExt,    // Order book subscriptions
			klineExt, // Kline subscriptions
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created spot realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
