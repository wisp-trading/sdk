package types

import (
	market2 "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// MarketStore handles perpetual market data storage.
// Embeds base MarketStore and all perp-specific extensions.
type MarketStore interface {
	market2.MarketStore
	market2.OrderBookStoreExtension
	market2.KlineStoreExtension
	FundingRateStoreExtension
}

// Perp-specific data keys
const (
	DataKeyFundingRates      market2.DataKey = "funding_rates"
	DataKeyHistoricalFunding market2.DataKey = "historical_funding"
)

// Perp-specific type aliases
type FundingRateMap map[connector.ExchangeName]perpConn.FundingRate
type HistoricalFundingMap map[connector.ExchangeName][]perpConn.HistoricalFundingRate

// FundingRateStoreExtension is the perp-specific store extension for funding rate data.
type FundingRateStoreExtension interface {
	market2.StoreExtension

	// Current funding rates
	UpdateFundingRate(asset portfolio.Pair, exchange connector.ExchangeName, rate perpConn.FundingRate)
	UpdateFundingRates(exchange connector.ExchangeName, rates map[portfolio.Pair]perpConn.FundingRate)
	GetFundingRate(asset portfolio.Pair, exchange connector.ExchangeName) *perpConn.FundingRate
	GetFundingRatesForAsset(asset portfolio.Pair) map[connector.ExchangeName]perpConn.FundingRate
	GetAllAssetsWithFundingRates() []portfolio.Pair

	// Historical funding rates
	UpdateHistoricalFundingRates(asset portfolio.Pair, exchange connector.ExchangeName, rates []perpConn.HistoricalFundingRate)
	GetHistoricalFundingRates(asset portfolio.Pair, exchange connector.ExchangeName) []perpConn.HistoricalFundingRate
	GetHistoricalFundingRatesForAsset(asset portfolio.Pair) map[connector.ExchangeName][]perpConn.HistoricalFundingRate
}
