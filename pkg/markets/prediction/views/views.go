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

// GetMarketViews returns all prediction markets currently on the watchlist,
func (v *predictionViews) GetMarketViews() []monitoring.PredictionMarketView {
	allMarkets := v.watchlist.GetAllMarkets()
	result := make([]monitoring.PredictionMarketView, 0)

	for exchange, markets := range allMarkets {
		for _, market := range markets {
			outcomes := make([]monitoring.PredictionOutcomeView, 0, len(market.Outcomes))
			for _, o := range market.Outcomes {
				outcomes = append(outcomes, monitoring.PredictionOutcomeView{
					OutcomeID: string(o.OutcomeID),
					Name:      o.Pair.Outcome(),
				})
			}
			result = append(result, monitoring.PredictionMarketView{
				Exchange: string(exchange),
				MarketID: string(market.MarketID),
				Slug:     market.Slug,
				Outcomes: outcomes,
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

// GetMarketOrderBooks returns all outcome order books for a given market.
func (v *predictionViews) GetMarketOrderBooks(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
) predictionStore.OutcomeOrderBookMap {
	return v.store.GetMarketOrderBooks(exchange, marketID)
}
