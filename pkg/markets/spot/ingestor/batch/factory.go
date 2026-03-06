package batch

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/ingestor/batch"
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type factory struct {
	connectorRegistry registry.ConnectorRegistry
	watchlist         spotTypes.SpotWatchlist
	store             spotTypes.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	watchlist spotTypes.SpotWatchlist,
	store spotTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) spotTypes.SpotBatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

// CreateIngestors creates one batch ingestor per registered spot connector
func (f *factory) CreateIngestors() []batchTypes.BatchIngestor {
	spotConnectors := f.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().Build())

	ingestorList := make([]batchTypes.BatchIngestor, 0, len(spotConnectors))

	for _, conn := range spotConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		marketDataReader, ok := conn.(connector.MarketDataReader)
		if !ok {
			f.logger.Warn("Spot connector %s does not implement MarketDataReader, skipping batch ingestor", exchangeName)
			continue
		}

		klineExt := batch.NewKlineExtension(
			marketDataReader,
			f.store,
			f.logger,
			[]string{"1m", "5m", "15m", "1h"},
			map[string]int{
				"1m":  500,
				"5m":  300,
				"15m": 200,
				"1h":  168,
			},
		)
		priceExt := batch.NewPriceExtension(marketDataReader, f.store, f.logger)
		orderbookExt := batch.NewOrderBookExtension(marketDataReader, f.store, f.logger, 20)

		ingestor := batch.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypeSpot,
			f.watchlist,
			f.store,
			f.timeProvider,
			f.logger,
			klineExt,
			priceExt,
			orderbookExt,
		)

		ingestorList = append(ingestorList, ingestor)

		f.logger.Info("Created spot batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
