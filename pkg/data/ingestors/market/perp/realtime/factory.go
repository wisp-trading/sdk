package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/real_time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	perpStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered perp WebSocket connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	marketWatchlist   data.MarketWatchlist
	store             perpStore.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	marketWatchlist data.MarketWatchlist,
	store perpStore.MarketStore,
	logger logging.ApplicationLogger,
) realtime.RealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		marketWatchlist:   marketWatchlist,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered perp WebSocket connector
func (f *factory) CreateIngestors() []realtime.RealtimeIngestor {
	perpWSConnectors := f.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtime.RealtimeIngestor, 0, len(perpWSConnectors))

	for _, wsConn := range perpWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		// Create extensions for perp markets
		obExt := real_time.NewOrderBookExtension(f.store, f.logger)
		klineExt := real_time.NewKlineExtension(f.store, f.logger, []string{"1m", "5m", "15m", "1h"})
		fundingExt := NewFundingRateExtension(f.store, f.logger)

		// Base ingestor + perp extensions
		ingestor := real_time.NewRealtimeIngestor(
			wsConn, // Perp WebSocket connector
			exchangeName,
			connector.MarketTypePerp,
			f.marketWatchlist,
			f.logger,
			obExt,      // Order book subscriptions
			klineExt,   // Kline subscriptions
			fundingExt, // Funding rate subscriptions
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created perp realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
