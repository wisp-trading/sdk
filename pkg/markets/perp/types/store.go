package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// MarketStore handles perpetual market data storage.
// Embeds base MarketStore and all perp-specific extensions.
type MarketStore interface {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	FundingRateStoreExtension
	PerpPositionsStoreExtension
}

// FundingRateStoreExtension is the perp-specific store extension for funding rate data.
type FundingRateStoreExtension interface {
	market.StoreExtension

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

// PerpPositionsStoreExtension stores live perp positions as reported by the exchange.
// Written by the realtime ingestor; read by the SDK and PNL calculator.
type PerpPositionsStoreExtension interface {
	market.StoreExtension

	// UpsertPosition inserts or replaces a position for the given exchange + pair.
	UpsertPosition(position perpConn.Position)

	// RemovePosition removes an existing position (called when size reaches zero).
	RemovePosition(exchange connector.ExchangeName, pair portfolio.Pair)

	// GetPositions returns all known live positions.
	GetPositions() []perpConn.Position

	// GetPosition returns the position for a specific exchange + pair, or nil.
	GetPosition(exchange connector.ExchangeName, pair portfolio.Pair) *perpConn.Position

	// QueryPositions returns positions filtered by exchange and/or pair.
	QueryPositions(q market.ActivityQuery) []perpConn.Position
}
