package types

import (
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type Predict interface {
	GetMarketBySlug(slug string, exchange connector.ExchangeName) (predictionconnector.Market, error)
	GetRecurringMarketBySlug(slug string, recurrenceInterval predictionconnector.RecurrenceInterval, exchange connector.ExchangeName) (predictionconnector.Market, error)
	WatchMarket(exchange connector.ExchangeName, market predictionconnector.Market)

	Markets() []predictionconnector.Market

	Orderbook(
		exchange connector.ExchangeName,
		market predictionconnector.Market,
		outcome predictionconnector.Outcome,
	) (*connector.OrderBook, error)

	// Balance returns the current balance for an asset on an exchange.
	Balance(exchange connector.ExchangeName, asset portfolio.Asset) (numerical.Decimal, bool)

	// Positions returns all orders recorded for the given strategy.
	Positions(strategyName strategy.StrategyName) []PredictionOrder

	// Log returns the trading logger for strategy-specific logging.
	// Use for recording trading decisions and strategy events.
	Log() logging.TradingLogger

	// PredictionSignal creates a new signal builder for prediction market trading signals.
	// Example: k.PredictionSignal(strategyName).Buy(market, outcome, exchange, shares, maxPrice, expiry).Build()
	PredictionSignal(strategyName strategy.StrategyName) PredictionSignalBuilder

	GetTokensToRedeem(market predictionconnector.Market) ([]predictionconnector.Balance, error)

	// Redeem attempts to redeem winnings for a market. Returns an error if redemption fails.
	Redeem(market predictionconnector.Market) error
}
