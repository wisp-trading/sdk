package predict

import (
	"errors"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	store "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// Wisp is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
type predict struct {
	applicationLogger   logging.ApplicationLogger
	tradingLogger       logging.TradingLogger
	signal              types.SignalFactory
	store               store.MarketStore
	predictionWatchlist types.PredictionWatchlist
	connectorRegistry   registry.ConnectorRegistry
}

// NewPredict constructs a new Predict instance with the provided dependencies.
func NewPredict(
	applicationLogger logging.ApplicationLogger,
	tradingLogger logging.TradingLogger,
	signal types.SignalFactory,
	store store.MarketStore,
	predictionWatchlist types.PredictionWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) types.Predict {
	return &predict{
		applicationLogger:   applicationLogger,
		tradingLogger:       tradingLogger,
		signal:              signal,
		predictionWatchlist: predictionWatchlist,
		connectorRegistry:   connectorRegistry,
		store:               store,
	}
}

func (p predict) GetMarketBySlug(slug string, exchange connector.ExchangeName) (predictionconnector.Market, error) {
	marketConnector, exists := p.connectorRegistry.Prediction(exchange)

	if !exists {
		return predictionconnector.Market{}, errors.New("connector not found for exchange: " + string(exchange))
	}

	market, err := marketConnector.GetMarket(slug)

	if err != nil {
		return predictionconnector.Market{}, errors.New("market not found for slug: " + slug + " on exchange: " + string(exchange))
	}

	return market, nil
}

func (p predict) GetRecurringMarketBySlug(
	slug string,
	recurrenceInterval predictionconnector.RecurrenceInterval,
	exchange connector.ExchangeName,
) (predictionconnector.Market, error) {
	marketConnector, exists := p.connectorRegistry.Prediction(exchange)

	if !exists {
		return predictionconnector.Market{}, errors.New("connector not found for exchange: " + string(exchange))
	}

	market, err := marketConnector.GetRecurringMarket(slug, recurrenceInterval)

	if err != nil {
		return predictionconnector.Market{}, errors.New("market not found for slug: " + slug + " on exchange: " + string(exchange))
	}

	return market, nil
}

func (p predict) WatchMarket(exchange connector.ExchangeName, market predictionconnector.Market) {
	p.predictionWatchlist.RequireMarket(exchange, market)
}

func (p predict) Markets() []predictionconnector.Market {
	//TODO implement me
	panic("implement me")
}

func (p predict) Orderbooks(market predictionconnector.Market) (map[predictionconnector.Outcome]predictionconnector.OrderBook, error) {
	//TODO implement me
	panic("implement me")
}

func (p predict) Orderbook(exchange connector.ExchangeName, market predictionconnector.Market, outcome predictionconnector.Outcome) (*connector.OrderBook, error) {
	book := p.store.GetOrderBook(exchange, market.MarketID, outcome.OutcomeID)

	if book == nil {
		return nil, errors.New("order book not found for outcome: " + string(outcome.OutcomeID) + " on exchange: " + string(exchange))
	}

	return book, nil
}

func (p predict) Log() logging.TradingLogger {
	return p.tradingLogger
}

// PredictionSignal creates a new signal builder for prediction market trading signals.
func (p predict) PredictionSignal(strategyName strategy.StrategyName) types.PredictionSignalBuilder {
	return p.signal.NewPrediction(strategyName)
}
