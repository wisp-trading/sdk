package types

import (
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	predictionStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
)

// PredictionViews owns all monitoring view logic for prediction markets.
// The monitoring ViewRegistry delegates to this interface — it does not implement
// prediction-specific logic itself.
type PredictionViews interface {
	// GetMarketViews returns all prediction markets currently being watched,
	// structured as PredictionMarketView entries with their full outcome lists.
	// Driven live from the prediction watchlist — never a stale snapshot.
	GetMarketViews() []monitoring.PredictionMarketView

	// GetOrderBook returns the order book for a specific outcome on a prediction market.
	GetOrderBook(
		exchange connector.ExchangeName,
		marketID predictionconnector.MarketID,
		outcomeID predictionconnector.OutcomeID,
	) *connector.OrderBook

	// GetMarketOrderBooks returns all outcome order books for a given market.
	GetMarketOrderBooks(
		exchange connector.ExchangeName,
		marketID predictionconnector.MarketID,
	) predictionStore.OutcomeOrderBookMap
}
