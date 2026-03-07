package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// Perp is the domain-scoped context object for perpetual futures strategies.
// Injected via wisp.Perp() — owns watchlist management, market data, positions,
// and signal creation for the perp domain.
type Perp interface {
	// WatchPair registers a pair on the perp watchlist so data ingestors begin
	// collecting orderbook, kline, and funding rate data for it.
	WatchPair(exchange connector.ExchangeName, pair portfolio.Pair)

	// UnwatchPair removes a pair from the perp watchlist.
	UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair)

	// FundingRate returns the latest funding rate for a pair on a specific exchange.
	FundingRate(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.FundingRate, bool)

	// FundingRates returns the funding rate across all exchanges for a pair.
	FundingRates(pair portfolio.Pair) map[connector.ExchangeName]perpConn.FundingRate

	// Position returns the current open position for a pair on an exchange, if any.
	Position(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.Position, bool)

	// Positions returns live perp positions from the store.
	// Optionally filter by exchange and/or pair: wisp.Perp().Positions()
	//   wisp.Perp().Positions(market.ActivityQuery{Exchange: &exchange})
	//   wisp.Perp().Positions(market.ActivityQuery{Pair: &pair})
	Positions(q ...market.ActivityQuery) []perpConn.Position

	// OrderBook returns the latest order book for a pair on a specific exchange.
	OrderBook(exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, bool)

	// Klines returns historical kline data for a pair on a specific exchange.
	Klines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline

	// Signal creates a new perp signal builder for the given strategy.
	Signal(strategyName strategy.StrategyName) strategy.PerpSignalBuilder

	// Log returns the trading logger for strategy-specific logging.
	Log() logging.TradingLogger

	// Trades returns all trades executed in the perp domain.
	// Optionally filter by exchange and/or pair.
	Trades(q ...market.ActivityQuery) []connector.Trade

	// PNL returns profit and loss calculations for this perp instance.
	PNL() PerpPNL
}
