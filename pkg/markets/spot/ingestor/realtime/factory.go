package realtime

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	spotStore "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered spot WebSocket connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	marketWatchlist   types.MarketWatchlist
	store             spotStore.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	marketWatchlist types.MarketWatchlist,
	store spotStore.MarketStore,
	logger logging.ApplicationLogger,
) realtimeTypes.RealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		marketWatchlist:   marketWatchlist,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered spot WebSocket connector
func (f *factory) CreateIngestors() []realtimeTypes.RealtimeIngestor {
	spotWSConnectors := f.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtimeTypes.RealtimeIngestor, 0, len(spotWSConnectors))

	for _, wsConn := range spotWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		// Create extensions for spot markets
		obExt := realtime.NewOrderBookExtension(f.store, f.logger)
		klineExt := realtime.NewKlineExtension(f.store, f.logger, []string{"1m", "5m", "15m", "1h"})

		ingestor := realtime.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypeSpot,
			f.marketWatchlist,
			f.logger,
			obExt,    // Order book subscriptions
			klineExt, // Kline subscriptions
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created spot realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
