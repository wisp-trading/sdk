package views

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	predictionStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type predictionViews struct {
	watchlist         types.PredictionWatchlist
	store             predictionStore.MarketStore
	connectorRegistry registry.ConnectorRegistry
}

// NewPredictionViews constructs a PredictionViews implementation.
func NewPredictionViews(
	watchlist types.PredictionWatchlist,
	store predictionStore.MarketStore,
	connectorRegistry registry.ConnectorRegistry,
) types.PredictionViews {
	return &predictionViews{
		watchlist:         watchlist,
		store:             store,
		connectorRegistry: connectorRegistry,
	}
}

// GetAvailableMarkets returns all prediction markets currently registered on the watchlist,
// formatted as AssetExchange entries for the monitoring layer.
func (v *predictionViews) GetAvailableMarkets() []monitoring.AssetExchange {
	markets := v.watchlist.GetAllMarkets()
	result := make([]monitoring.AssetExchange, 0, len(markets))

	for exchange, marketList := range markets {
		for _, market := range marketList {
			result = append(result, monitoring.AssetExchange{
				Asset:    string(market.MarketID),
				Exchange: string(exchange),
			})
		}
	}

	return result
}

// GetOrderBook returns the order book for a specific outcome on a prediction market.
func (v *predictionViews) GetOrderBook(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
	outcomeID predictionconnector.OutcomeID,
) *connector.OrderBook {
	return v.store.GetOrderBook(exchange, marketID, outcomeID)
}

// GetMarketOrderBooks returns all known outcome order books for a market.
func (v *predictionViews) GetMarketOrderBooks(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
) predictionStore.OutcomeOrderBookMap {
	return v.store.GetMarketOrderBooks(exchange, marketID)
}
