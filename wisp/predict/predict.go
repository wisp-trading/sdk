package predict

import (
	"errors"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	predictTypes "github.com/wisp-trading/sdk/pkg/types/wisp/predict"
)

// Wisp is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
type predict struct {
	applicationLogger logging.ApplicationLogger
	tradingLogger     logging.TradingLogger
	signal            strategy.SignalFactory
	connectorRegistry registry.ConnectorRegistry
}

// NewPredict constructs a new Predict instance with the provided dependencies.
func NewPredict(
	applicationLogger logging.ApplicationLogger,
	tradingLogger logging.TradingLogger,
	signal strategy.SignalFactory,
	connectorRegistry registry.ConnectorRegistry,
) predictTypes.Predict {
	return &predict{
		applicationLogger: applicationLogger,
		tradingLogger:     tradingLogger,
		signal:            signal,
		connectorRegistry: connectorRegistry,
	}
}

func (p predict) GetMarketBySlug(slug string, exchange connector.ExchangeName) (prediction.Market, error) {
	marketConnector, exists := p.connectorRegistry.Prediction(exchange)

	if !exists {
		return prediction.Market{}, errors.New("connector not found for exchange: " + string(exchange))
	}

	market, err := marketConnector.GetMarket(slug)

	if err != nil {
		return prediction.Market{}, errors.New("market not found for slug: " + slug + " on exchange: " + string(exchange))
	}

	return market, nil
}

func (p predict) GetRecurringMarketBySlug(
	slug string,
	recurrenceInterval prediction.RecurrenceInterval,
	exchange connector.ExchangeName,
) (prediction.Market, error) {
	marketConnector, exists := p.connectorRegistry.Prediction(exchange)

	if !exists {
		return prediction.Market{}, errors.New("connector not found for exchange: " + string(exchange))
	}

	market, err := marketConnector.GetRecurringMarket(slug, recurrenceInterval)

	if err != nil {
		return prediction.Market{}, errors.New("market not found for slug: " + slug + " on exchange: " + string(exchange))
	}

	return market, nil
}

func (p predict) WatchMarket(market prediction.Market, exchange *connector.ExchangeName) error {
	if exchange != nil {
		marketConnector, exists := p.connectorRegistry.PredictionWebSocket(*exchange)

		if !exists {
			return errors.New("connector not found for exchange: " + string(*exchange))
		}

		return marketConnector.SubscribeOrderBook(market)
	}

	connectors := p.connectorRegistry.FilterPrediction(registry.NewFilter().ReadyOnly().WebSocketOnly().Build())
	for _, conn := range connectors {
		marketConnector, exists := p.connectorRegistry.PredictionWebSocket(conn.GetConnectorInfo().Name)

		if !exists {
			p.applicationLogger.Warn("WebSocket connector not found for exchange: " + string(conn.GetConnectorInfo().Name))
			continue
		}

		err := marketConnector.SubscribeOrderBook(market)
		if err != nil {
			p.applicationLogger.Warn("Failed to subscribe to market on connector: " + err.Error())
			continue
		}
	}

	return nil
}

func (p predict) Markets() []prediction.Market {
	//TODO implement me
	panic("implement me")
}

func (p predict) Orderbooks(market prediction.Market) (map[prediction.Outcome]prediction.OrderBook, error) {
	//TODO implement me
	panic("implement me")
}

func (p predict) Orderbook(market prediction.Market, outcome prediction.Outcome) (*prediction.OrderBook, error) {
	//TODO implement me
	panic("implement me")
}

func (p predict) Log() logging.TradingLogger {
	return p.tradingLogger
}
