package batch

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	perpConn "github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/batch"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// FundingRateExtension collects perp-specific funding rate data
type FundingRateExtension struct {
	store  perpStore.MarketStore
	logger logging.ApplicationLogger
}

func NewFundingRateExtension(store perpStore.MarketStore, logger logging.ApplicationLogger) batch.CollectionExtension {
	return &FundingRateExtension{
		store:  store,
		logger: logger,
	}
}

func (f *FundingRateExtension) Collect(conn connector.Connector, exchangeName connector.ExchangeName, assets []portfolio.Asset) {
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

var _ batch.CollectionExtension = (*FundingRateExtension)(nil)
