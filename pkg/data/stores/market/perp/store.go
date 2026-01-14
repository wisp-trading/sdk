package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market/extensions"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	perpTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// perpStore wraps the base store and adds perp-specific methods
type perpStore struct {
	market.MarketStore
	fundingExt *extensions.FundingRateExtension
}

// NewStore creates a new perp market store with funding rate extension
func NewStore(timeProvider temporal.TimeProvider) perpTypes.MarketStore {
	fundingExt := extensions.NewFundingRateExtension()

	baseStore := market.NewStore(timeProvider, fundingExt)

	return &perpStore{
		MarketStore: baseStore,
		fundingExt:  fundingExt,
	}
}

// ===== Perp-specific methods (funding rates) =====

func (ps *perpStore) UpdateFundingRate(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
	rate connector.FundingRate,
) {
	ps.fundingExt.UpdateFundingRate(asset, exchange, rate)
}

func (ps *perpStore) UpdateFundingRates(
	exchange connector.ExchangeName,
	rates map[portfolio.Asset]connector.FundingRate,
) {
	ps.fundingExt.UpdateFundingRates(exchange, rates)
}

func (ps *perpStore) GetFundingRate(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
) *connector.FundingRate {
	return ps.fundingExt.GetFundingRate(asset, exchange)
}

func (ps *perpStore) GetFundingRatesForAsset(
	asset portfolio.Asset,
) perpTypes.FundingRateMap {
	return ps.fundingExt.GetFundingRatesForAsset(asset)
}

func (ps *perpStore) GetAllAssetsWithFundingRates() []portfolio.Asset {
	return ps.fundingExt.GetAllAssetsWithFundingRates()
}

func (ps *perpStore) UpdateHistoricalFundingRates(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
	rates []connector.HistoricalFundingRate,
) {
	ps.fundingExt.UpdateHistoricalFundingRates(asset, exchange, rates)
}

func (ps *perpStore) GetHistoricalFundingRates(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
) []connector.HistoricalFundingRate {
	return ps.fundingExt.GetHistoricalFundingRates(asset, exchange)
}

func (ps *perpStore) GetHistoricalFundingRatesForAsset(
	asset portfolio.Asset,
) perpTypes.HistoricalFundingMap {
	return ps.fundingExt.GetHistoricalFundingRatesForAsset(asset)
}
