package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// MarketStore handles perpetual market data storage
// Embeds base.MarketStore and adds perp-specific methods
type MarketStore interface {
	market.MarketStore

	// Funding rates
	UpdateFundingRate(asset portfolio.Asset, exchange connector.ExchangeName, rate perp.FundingRate)
	UpdateFundingRates(exchange connector.ExchangeName, rates map[portfolio.Asset]perp.FundingRate)
	GetFundingRate(asset portfolio.Asset, exchange connector.ExchangeName) *perp.FundingRate
	GetFundingRatesForAsset(asset portfolio.Asset) FundingRateMap
	GetAllAssetsWithFundingRates() []portfolio.Asset

	// Historical funding rates
	UpdateHistoricalFundingRates(asset portfolio.Asset, exchange connector.ExchangeName, rates []perp.HistoricalFundingRate)
	GetHistoricalFundingRates(asset portfolio.Asset, exchange connector.ExchangeName) []perp.HistoricalFundingRate
	GetHistoricalFundingRatesForAsset(asset portfolio.Asset) HistoricalFundingMap
}

// Perp-specific data keys (extend base.DataKey)
const (
	DataKeyFundingRates      market.DataKey = "funding_rates"
	DataKeyHistoricalFunding market.DataKey = "historical_funding"
	DataKeyContract          market.DataKey = "contract"
)

// Perp-specific type aliases (from old data_store.go)
type FundingRateMap map[connector.ExchangeName]perp.FundingRate
type HistoricalFundingMap map[connector.ExchangeName][]perp.HistoricalFundingRate
