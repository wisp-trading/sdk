package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// MarketMap is Exchange -> Market
type MarketMap map[connector.ExchangeName]map[predictionconnector.MarketID]predictionconnector.Market

// MarketStoreExtension provides market metadata storage for discovery and screening
type MarketStoreExtension interface {
	market.StoreExtension

	// UpdateMarkets stores or updates a batch of markets for an exchange
	UpdateMarkets(exchange connector.ExchangeName, markets []predictionconnector.Market)

	// GetMarket retrieves a single market by ID
	GetMarket(exchange connector.ExchangeName, marketID predictionconnector.MarketID) *predictionconnector.Market

	// GetMarkets retrieves all markets for an exchange
	GetMarkets(exchange connector.ExchangeName) []predictionconnector.Market

	// GetAllMarkets retrieves all markets across all exchanges
	GetAllMarkets() MarketMap

	// ClearMarkets removes all markets for an exchange
	ClearMarkets(exchange connector.ExchangeName)
}
