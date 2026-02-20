package batch

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	perpStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// fundingRateExtension collects perp-specific funding rate data
type fundingRateExtension struct {
	store  perpStore.MarketStore
	logger logging.ApplicationLogger
}

func NewFundingRateExtension(store perpStore.MarketStore, logger logging.ApplicationLogger) batch.CollectionExtension {
	return &fundingRateExtension{
		store:  store,
		logger: logger,
	}
}

func (f *fundingRateExtension) Collect(conn connector.Connector, exchangeName connector.ExchangeName, assets []portfolio.Pair) {
	pc, ok := conn.(perpConn.Connector)
	if !ok {
		f.logger.Debug("Connector %s does not support perp operations", exchangeName)
		return
	}

	// Fetch current funding rates
	rates, err := pc.FetchCurrentFundingRates()
	if err != nil {
		f.logger.Error("Failed to fetch funding rates from %s: %v", exchangeName, err)
		return
	}

	// Update all funding rates from this connector
	f.store.UpdateFundingRates(exchangeName, rates)

	for asset, rate := range rates {
		f.logger.Debug("Updated funding rate for %s on %s = %s",
			asset.Symbol(), exchangeName, rate.CurrentRate.String())
	}
}

var _ batch.CollectionExtension = (*fundingRateExtension)(nil)
