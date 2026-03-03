package predict

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type Predict interface {
	GetMarketBySlug(slug string, exchange connector.ExchangeName) (prediction.Market, error)
	GetRecurringMarketBySlug(slug string, recurrenceInterval prediction.RecurrenceInterval, exchange connector.ExchangeName) (prediction.Market, error)
	WatchMarket(exchange connector.ExchangeName, market prediction.Market)

	Markets() []prediction.Market
	//Orderbooks(market prediction.Market) (map[prediction.Outcome]prediction.OrderBook, error)
	Orderbook(
		exchange connector.ExchangeName,
		market prediction.Market,
		outcome prediction.Outcome,
	) (*connector.OrderBook, error)

	// Log returns the trading logger for strategy-specific logging.
	// Use for recording trading decisions and strategy events.
	Log() logging.TradingLogger

	// PredictionSignal creates a new signal builder for prediction market trading signals.
	// Example: k.PredictionSignal(strategyName).Buy(market, outcome, exchange, shares, maxPrice, expiry).Build()
	PredictionSignal(strategyName strategy.StrategyName) strategy.PredictionSignalBuilder
}
