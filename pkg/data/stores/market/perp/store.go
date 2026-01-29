package perp

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/perp/extensions"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// perpStore wraps the base store and adds perp-specific methods
type perpStore struct {
	marketTypes.MarketStore
	fundingExt *extensions.FundingRateExtension
}

// NewStore creates a new perp market store with funding rate extension
func NewStore(timeProvider temporal.TimeProvider) perpTypes.MarketStore {
	fundingExt := extensions.NewFundingRateExtension()

	baseStore := store.NewStore(timeProvider, fundingExt)

	return &perpStore{
		MarketStore: baseStore,
		fundingExt:  fundingExt,
	}
}

// ===== Perp-specific methods (funding rates) =====

func (ps *perpStore) UpdateFundingRate(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
	rate perp.FundingRate,
) {
	ps.fundingExt.UpdateFundingRate(asset, exchange, rate)
}

func (ps *perpStore) UpdateFundingRates(
	exchange connector.ExchangeName,
	rates map[portfolio.Asset]perp.FundingRate,
) {
	ps.fundingExt.UpdateFundingRates(exchange, rates)
}

func (ps *perpStore) GetFundingRate(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
) *perp.FundingRate {
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
	rates []perp.HistoricalFundingRate,
) {
	ps.fundingExt.UpdateHistoricalFundingRates(asset, exchange, rates)
}

func (ps *perpStore) GetHistoricalFundingRates(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
) []perp.HistoricalFundingRate {
	return ps.fundingExt.GetHistoricalFundingRates(asset, exchange)
}

func (ps *perpStore) GetHistoricalFundingRatesForAsset(
	asset portfolio.Asset,
) perpTypes.HistoricalFundingMap {
	return ps.fundingExt.GetHistoricalFundingRatesForAsset(asset)
}

// MarketType returns the spot market type
func (s *perpStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypePerp
}
