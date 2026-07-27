// Package facade provides shared pair-market methods for spot/perp domain facades.
// Domain packages embed Core and add Signal/Emit/PNL (and perp funding/positions).
package facade

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PairStore is the store surface used by pair domain facades for market data + trades/orders.
type PairStore interface {
	GetPairPrice(pair portfolio.Pair, exchange connector.ExchangeName) *connector.Price
	GetPairPrices(pair portfolio.Pair) market.PriceMap
	GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetKlines(pair portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
	GetAllTrades() []connector.Trade
	QueryTrades(q market.ActivityQuery) []connector.Trade
	GetOrders() []connector.Order
	QueryOrders(q market.ActivityQuery) []connector.Order
}

// Core holds shared deps and implements pair watchlist + market-data reads.
type Core struct {
	Logger       logging.TradingLogger
	Watchlist    baseTypes.MarketWatchlist
	Store        PairStore
	TimeProvider temporal.TimeProvider
	Router       execution.SignalRouter
}

// WatchPair registers a pair so ingestors start collecting data.
func (c *Core) WatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	c.Watchlist.RequirePair(exchange, pair)
}

// UnwatchPair removes a pair from the watchlist.
func (c *Core) UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	c.Watchlist.ReleasePair(exchange, pair)
}

// Price returns the current price for a pair on an exchange.
func (c *Core) Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool) {
	price := c.Store.GetPairPrice(pair, exchange)
	if price == nil {
		return numerical.Zero(), false
	}
	return price.Price, true
}

// Prices returns prices across exchanges for a pair.
func (c *Core) Prices(pair portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	priceMap := c.Store.GetPairPrices(pair)
	out := make(map[connector.ExchangeName]numerical.Decimal, len(priceMap))
	for exchange, p := range priceMap {
		out[exchange] = p.Price
	}
	return out
}

// OrderBook returns the latest order book, or false if missing.
func (c *Core) OrderBook(exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, bool) {
	ob := c.Store.GetOrderBook(pair, exchange)
	if ob == nil {
		return nil, false
	}
	return ob, true
}

// Klines returns historical klines from the store.
func (c *Core) Klines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return c.Store.GetKlines(pair, exchange, interval, limit)
}

// Trades returns trades, optionally filtered.
func (c *Core) Trades(q ...market.ActivityQuery) []connector.Trade {
	if len(q) > 0 {
		return c.Store.QueryTrades(q[0])
	}
	return c.Store.GetAllTrades()
}

// Orders returns placed orders, optionally filtered (spot Positions surface).
func (c *Core) Orders(q ...market.ActivityQuery) []connector.Order {
	if len(q) > 0 {
		return c.Store.QueryOrders(q[0])
	}
	return c.Store.GetOrders()
}

// Log returns the trading logger.
func (c *Core) Log() logging.TradingLogger {
	return c.Logger
}

// Emit routes a signal through the shared executor router.
func (c *Core) Emit(signal strategy.Signal) execution.ExecutionCallback {
	return execution.Dispatch(c.Router, signal)
}
