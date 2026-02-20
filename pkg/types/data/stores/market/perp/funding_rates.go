package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// Perp-specific data keys
const (
	DataKeyFundingRates      marketTypes.DataKey = "funding_rates"
	DataKeyHistoricalFunding marketTypes.DataKey = "historical_funding"
)

// Perp-specific type aliases
type FundingRateMap map[connector.ExchangeName]perp.FundingRate
type HistoricalFundingMap map[connector.ExchangeName][]perp.HistoricalFundingRate

type FundingRateStoreExtension interface {
	marketTypes.StoreExtension

	// Current funding rates
	UpdateFundingRate(asset portfolio.Pair, exchange connector.ExchangeName, rate perp.FundingRate)
	UpdateFundingRates(exchange connector.ExchangeName, rates map[portfolio.Pair]perp.FundingRate)
	GetFundingRate(asset portfolio.Pair, exchange connector.ExchangeName) *perp.FundingRate
	GetFundingRatesForAsset(asset portfolio.Pair) map[connector.ExchangeName]perp.FundingRate
	GetAllAssetsWithFundingRates() []portfolio.Pair

	// Historical funding rates
	UpdateHistoricalFundingRates(asset portfolio.Pair, exchange connector.ExchangeName, rates []perp.HistoricalFundingRate)
	GetHistoricalFundingRates(asset portfolio.Pair, exchange connector.ExchangeName) []perp.HistoricalFundingRate
	GetHistoricalFundingRatesForAsset(asset portfolio.Pair) map[connector.ExchangeName][]perp.HistoricalFundingRate
}
