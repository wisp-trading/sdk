package realtime

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/real_time"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// factory creates realtime ingestors for all registered perp WebSocket connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	watchlist         perpTypes.PerpWatchlist
	store             perpTypes.MarketStore
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	watchlist perpTypes.PerpWatchlist,
	store perpTypes.MarketStore,
	logger logging.ApplicationLogger,
) realtimeTypes.RealtimeIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		store:             store,
		logger:            logger,
	}
}

// CreateIngestors creates one realtime ingestor per registered perp WebSocket connector
func (f *factory) CreateIngestors() []realtimeTypes.RealtimeIngestor {
	perpWSConnectors := f.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())

	realtimeIngestors := make([]realtimeTypes.RealtimeIngestor, 0, len(perpWSConnectors))

	for _, wsConn := range perpWSConnectors {
		info := wsConn.GetConnectorInfo()
		exchangeName := info.Name

		// Create extensions for perp markets
		obExt := real_time.NewOrderBookExtension(f.store, f.logger)
		klineExt := real_time.NewKlineExtension(f.store, f.logger, []string{"1m", "5m", "15m", "1h"})
		fundingExt := NewFundingRateExtension(f.store, f.logger)

		// Adapt PerpWatchlist to the MarketWatchlist interface expected by the base ingestor
		watchlistAdapter := newWatchlistAdapter(f.watchlist)

		// Base ingestor + perp extensions
		ingestor := real_time.NewRealtimeIngestor(
			wsConn,
			exchangeName,
			connector.MarketTypePerp,
			watchlistAdapter,
			f.logger,
			obExt,
			klineExt,
			fundingExt,
		)

		realtimeIngestors = append(realtimeIngestors, ingestor)

		f.logger.Info("Created perp realtime ingestor for %s", exchangeName)
	}

	return realtimeIngestors
}
