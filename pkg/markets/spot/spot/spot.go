package spot

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	spotSignal "github.com/wisp-trading/sdk/pkg/markets/spot/signal"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type spot struct {
	tradingLogger logging.TradingLogger
	watchlist     spotTypes.SpotWatchlist
	store         spotTypes.MarketStore
	timeProvider  temporal.TimeProvider
	pnl           spotTypes.SpotPNL
}

func NewSpot(
	tradingLogger logging.TradingLogger,
	watchlist spotTypes.SpotWatchlist,
	store spotTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	pnl spotTypes.SpotPNL,
) spotTypes.Spot {
	return &spot{
		tradingLogger: tradingLogger,
		watchlist:     watchlist,
		store:         store,
		timeProvider:  timeProvider,
		pnl:           pnl,
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

func (s *spot) Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool) {
	price := s.store.GetPairPrice(pair, exchange)
	if price == nil {
		return numerical.Zero(), false
	}
	return price.Price, true
}

func (s *spot) Prices(pair portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	priceMap := s.store.GetPairPrices(pair)
	out := make(map[connector.ExchangeName]numerical.Decimal, len(priceMap))
	for exchange, p := range priceMap {
		out[exchange] = p.Price
	}
	return out
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
func (s *spot) Signal(strategyName strategy.StrategyName) spotTypes.SpotSignalBuilder {
	return spotSignal.NewSpotBuilder(strategyName, s.timeProvider)
}

// Log returns the trading logger for strategy-specific logging.
func (s *spot) Log() logging.TradingLogger {
	return s.tradingLogger
}

func (s *spot) Trades(q ...market.ActivityQuery) []connector.Trade {
	if len(q) > 0 {
		return s.store.QueryTrades(q[0])
	}
	return s.store.GetAllTrades()
}

func (s *spot) Positions(q ...market.ActivityQuery) []connector.Order {
	if len(q) > 0 {
		return s.store.QueryOrders(q[0])
	}
	return s.store.GetOrders()
}

func (s *spot) PNL() spotTypes.SpotPNL {
	return s.pnl
}

var _ spotTypes.Spot = (*spot)(nil)
