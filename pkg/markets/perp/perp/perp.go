package perp

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/perp/signal"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type perp struct {
	tradingLogger logging.TradingLogger
	watchlist     perpTypes.PerpWatchlist
	store         perpTypes.MarketStore
	timeProvider  temporal.TimeProvider
	pnl           perpTypes.PerpPNL
}

func NewPerp(
	tradingLogger logging.TradingLogger,
	watchlist perpTypes.PerpWatchlist,
	store perpTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	pnl perpTypes.PerpPNL,
) perpTypes.Perp {
	return &perp{
		tradingLogger: tradingLogger,
		watchlist:     watchlist,
		store:         store,
		timeProvider:  timeProvider,
		pnl:           pnl,
	}
}

// WatchPair registers a pair on the perp watchlist, triggering data ingestion.
func (p *perp) WatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	p.watchlist.RequirePair(exchange, pair)
}

// UnwatchPair removes a pair from the perp watchlist.
func (p *perp) UnwatchPair(exchange connector.ExchangeName, pair portfolio.Pair) {
	p.watchlist.ReleasePair(exchange, pair)
}

// FundingRate returns the latest funding rate for a pair on a specific exchange.
func (p *perp) FundingRate(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.FundingRate, bool) {
	rate := p.store.GetFundingRate(pair, exchange)
	if rate == nil {
		return nil, false
	}
	return rate, true
}

// FundingRates returns the funding rate across all exchanges for a pair.
func (p *perp) FundingRates(pair portfolio.Pair) map[connector.ExchangeName]perpConn.FundingRate {
	return p.store.GetFundingRatesForAsset(pair)
}

// Price returns the current mark price for a pair on a specific exchange.
func (p *perp) Price(exchange connector.ExchangeName, pair portfolio.Pair) (numerical.Decimal, bool) {
	price := p.store.GetPairPrice(pair, exchange)
	if price == nil {
		return numerical.Zero(), false
	}
	return price.Price, true
}

// Prices returns the current mark price for a pair across all exchanges.
func (p *perp) Prices(pair portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	priceMap := p.store.GetPairPrices(pair)
	out := make(map[connector.ExchangeName]numerical.Decimal, len(priceMap))
	for exchange, price := range priceMap {
		out[exchange] = price.Price
	}
	return out
}

// Position returns a single live position for a specific exchange + pair from the store.
func (p *perp) Position(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.Position, bool) {
	pos := p.store.GetPosition(exchange, pair)
	if pos == nil {
		return nil, false
	}
	return pos, true
}

// Positions returns live positions from the store.
// Optionally filter by exchange and/or pair.
func (p *perp) Positions(q ...market.ActivityQuery) []perpConn.Position {
	if len(q) > 0 {
		return p.store.QueryPositions(q[0])
	}
	return p.store.GetPositions()
}

// OrderBook returns the latest order book for a pair on a specific exchange.
func (p *perp) OrderBook(exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, bool) {
	ob := p.store.GetOrderBook(pair, exchange)
	if ob == nil {
		return nil, false
	}
	return ob, true
}

func (p *perp) Klines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return p.store.GetKlines(pair, exchange, interval, limit)
}

// Signal creates a new perp signal builder for the given strategy.
func (p *perp) Signal(strategyName strategy.StrategyName) perpTypes.PerpSignalBuilder {
	return signal.NewPerpBuilder(strategyName, p.timeProvider)
}

// Log returns the trading logger for strategy-specific logging.
func (p *perp) Log() logging.TradingLogger {
	return p.tradingLogger
}

func (p *perp) Trades(q ...market.ActivityQuery) []connector.Trade {
	if len(q) > 0 {
		return p.store.QueryTrades(q[0])
	}
	return p.store.GetAllTrades()
}

// PNL returns the profit and loss calculator for the perp context.
func (p *perp) PNL() perpTypes.PerpPNL {
	return p.pnl
}

var _ perpTypes.Perp = (*perp)(nil)
