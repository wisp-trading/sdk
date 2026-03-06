package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Spot is the domain-scoped context object for spot market strategies.
// Injected via wisp.Spot() — owns watchlist management, market data,
// and signal creation for the spot domain.
type Spot interface {
	// WatchPair registers a pair on the spot watchlist so data ingestors begin
	// collecting orderbook and kline data for it.
	WatchPair(exchange connector.ExchangeName, pair portfolio.Pair)

	// UnwatchPair removes a pair from the spot watchlist.
	UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair)

	// Price returns the current price for a pair on a specific exchange.
	Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool)

	// OrderBook returns the latest order book for a pair on a specific exchange.
	OrderBook(exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, bool)

	// Klines returns historical kline data for a pair on a specific exchange.
	Klines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline

	// Signal creates a new spot signal builder for the given strategy.
	// Example: wisp.Spot().Signal(strategyName).Buy(pair, exchange, qty).Build()
	Signal(strategyName strategy.StrategyName) strategy.SpotSignalBuilder

	// Log returns the trading logger for strategy-specific logging.
	Log() logging.TradingLogger

	// PNL returns profit and loss calculations for this spot instance.
	PNL() wispActivity.SpotPNL
}
