package spot

import (
	dataTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spot is the concrete implementation of spotTypes.Spot.
// Injected into strategies via wisp.Spot().
type spot struct {
	tradingLogger logging.TradingLogger
	watchlist     dataTypes.MarketWatchlist
	store         spotTypes.MarketStore
	signal        strategy.SignalFactory
}

// NewSpot constructs the spot context object injected into strategies.
func NewSpot(
	tradingLogger logging.TradingLogger,
	watchlist dataTypes.MarketWatchlist,
	store spotTypes.MarketStore,
	signal strategy.SignalFactory,
) spotTypes.Spot {
	return &spot{
		tradingLogger: tradingLogger,
		watchlist:     watchlist,
		store:         store,
		signal:        signal,
	}
}

// WatchPair registers a pair on the spot watchlist, triggering data ingestion.
func (s *spot) WatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	s.watchlist.RequirePair(exchange, pair)
}

// UnwatchPair removes a pair from the spot watchlist.
func (s *spot) UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	s.watchlist.ReleasePair(exchange, pair)
}

// Price returns the current price for a pair on a specific exchange.
func (s *spot) Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool) {
	price := s.store.GetPairPrice(pair, exchange)
	if price == nil {
		return numerical.Zero(), false
	}
	return price.Price, true
}

// OrderBook returns the latest order book for a pair on a specific exchange.
func (s *spot) OrderBook(exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, bool) {
	ob := s.store.GetOrderBook(pair, exchange)
	if ob == nil {
		return nil, false
	}
	return ob, true
}

// Klines returns historical kline data for a pair on a specific exchange.
func (s *spot) Klines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return s.store.GetKlines(pair, exchange, interval, limit)
}

// Signal creates a new spot signal builder for the given strategy.
func (s *spot) Signal(strategyName strategy.StrategyName) strategy.SpotSignalBuilder {
	return s.signal.NewSpot(strategyName)
}

// Log returns the trading logger for strategy-specific logging.
func (s *spot) Log() logging.TradingLogger {
	return s.tradingLogger
}

var _ spotTypes.Spot = (*spot)(nil)
