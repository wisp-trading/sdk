package types

import (
	predictiontypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	predictionStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
)

// PredictionViews exposes monitoring data for prediction markets.
type PredictionViews interface {
	// GetAvailableMarkets returns all prediction markets currently being watched,
	// grouped by exchange.
	GetAvailableMarkets() []monitoring.AssetExchange

	// GetOrderBook returns the latest order book for a specific outcome of a prediction market.
	// Returns nil if no data is available.
	GetOrderBook(
		exchange connector.ExchangeName,
		marketID predictiontypes.MarketID,
		outcomeID predictiontypes.OutcomeID,
	) *connector.OrderBook

	// GetMarketOrderBooks returns all outcome order books for a given market.
	GetMarketOrderBooks(
		exchange connector.ExchangeName,
		marketID predictiontypes.MarketID,
	) predictionStore.OutcomeOrderBookMap
}
