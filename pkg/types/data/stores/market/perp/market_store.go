package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/common"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// MarketStore handles perpetual market data storage
// Embeds common.MarketStore and adds perp-specific methods
type MarketStore interface {
	common.MarketStore

	// Funding rates
	UpdateFundingRate(asset portfolio.Asset, exchange connector.ExchangeName, rate connector.FundingRate)
	UpdateFundingRates(exchange connector.ExchangeName, rates map[portfolio.Asset]connector.FundingRate)
	GetFundingRate(asset portfolio.Asset, exchange connector.ExchangeName) *connector.FundingRate
	GetFundingRatesForAsset(asset portfolio.Asset) FundingRateMap
	GetAllAssetsWithFundingRates() []portfolio.Asset

	// Historical funding rates
	UpdateHistoricalFundingRates(asset portfolio.Asset, exchange connector.ExchangeName, rates []connector.HistoricalFundingRate)
	GetHistoricalFundingRates(asset portfolio.Asset, exchange connector.ExchangeName) []connector.HistoricalFundingRate
	GetHistoricalFundingRatesForAsset(asset portfolio.Asset) HistoricalFundingMap
}

// Perp-specific data keys (extend common.DataKey)
const (
	DataKeyFundingRates      common.DataKey = "funding_rates"
	DataKeyHistoricalFunding common.DataKey = "historical_funding"
	DataKeyContract          common.DataKey = "contract"
)

// Perp-specific type aliases (from old data_store.go)
type FundingRateMap map[connector.ExchangeName]connector.FundingRate
type HistoricalFundingMap map[connector.ExchangeName][]connector.HistoricalFundingRate
