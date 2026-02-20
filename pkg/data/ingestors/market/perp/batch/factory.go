package batch

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	batchTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	perpStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// factory creates batch ingestors for all registered perp connectors
type factory struct {
	connectorRegistry registry.ConnectorRegistry
	assetRegistry     registry.PairRegistry
	store             perpStore.MarketStore
	timeProvider      temporal.TimeProvider
	logger            logging.ApplicationLogger
}

func NewFactory(
	connectorRegistry registry.ConnectorRegistry,
	assetRegistry registry.PairRegistry,
	store perpStore.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) batchTypes.BatchIngestorFactory {
	return &factory{
		connectorRegistry: connectorRegistry,
		assetRegistry:     assetRegistry,
		store:             store,
		timeProvider:      timeProvider,
		logger:            logger,
	}
}

// CreateIngestors creates one batch ingestor per registered perp connector
func (f *factory) CreateIngestors() []batchTypes.BatchIngestor {
	perpConnectors := f.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().Build())

	ingestorList := make([]batchTypes.BatchIngestor, 0, len(perpConnectors))

	for _, conn := range perpConnectors {
		info := conn.GetConnectorInfo()
		exchangeName := info.Name

		marketDataReader, ok := conn.(connector.MarketDataReader)
		if !ok {
			f.logger.Warn("Perp connector %s does not implement MarketDataReader, skipping batch ingestor", exchangeName)
			continue
		}

		// Create perp-specific extensions
		fundingExt := NewFundingRateExtension(f.store, f.logger)
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

		// Base ingestor + perp extensions
		ingestor := batch.NewBatchIngestor(
			conn,
			exchangeName,
			connector.MarketTypePerp,
			f.assetRegistry,
			f.store,
			f.timeProvider,
			f.logger,
			klineExt,
			priceExt,
			fundingExt,
			orderbookExt,
		)

		ingestorList = append(ingestorList, ingestor)

		f.logger.Info("Created perp batch ingestor for %s", exchangeName)
	}

	return ingestorList
}
