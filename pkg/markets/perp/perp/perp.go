package perp

import (
	storeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/markets/perp/signal"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// perp is the concrete implementation of perpTypes.Perp.
// Injected into strategies via wisp.Perp().
type perp struct {
	tradingLogger     logging.TradingLogger
	watchlist         perpTypes.PerpWatchlist
	store             perpTypes.MarketStore
	connectorRegistry registry.ConnectorRegistry
	timeProvider      temporal.TimeProvider
	pnl               perpTypes.PerpPNL
}

// NewPerp constructs the perp context object injected into strategies.
func NewPerp(
	tradingLogger logging.TradingLogger,
	watchlist perpTypes.PerpWatchlist,
	store perpTypes.MarketStore,
	connectorRegistry registry.ConnectorRegistry,
	timeProvider temporal.TimeProvider,
	trades storeTypes.Trades,
	pnl perpTypes.PerpPNL,
) perpTypes.Perp {
	return &perp{
		tradingLogger:     tradingLogger,
		watchlist:         watchlist,
		store:             store,
		connectorRegistry: connectorRegistry,
		timeProvider:      timeProvider,
		pnl:               pnl,
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

// Position returns the current open position for a pair on an exchange, if any.
// Fetches live from the connector — positions are not cached in the store.
func (p *perp) Position(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.Position, bool) {
	conn, exists := p.connectorRegistry.Connector(exchange)
	if !exists {
		return nil, false
	}

	pm, ok := conn.(perpConn.PositionManager)
	if !ok {
		return nil, false
	}

	positions, err := pm.GetPositions()
	if err != nil {
		return nil, false
	}

	for _, pos := range positions {
		if pos.Pair.Symbol() == pair.Symbol() {
			p := pos
			return &p, true
		}
	}

	return nil, false
}

// Positions returns all open positions across all exchanges for a strategy,
// by querying every ready perp connector.
func (p *perp) Positions() []perpConn.Position {
	perpConnectors := p.connectorRegistry.FilterPerp(
		registry.NewFilter().ReadyOnly().Build(),
	)

	var all []perpConn.Position
	for _, conn := range perpConnectors {
		pm, ok := conn.(perpConn.PositionManager)
		if !ok {
			continue
		}
		positions, err := pm.GetPositions()
		if err != nil {
			continue
		}
		all = append(all, positions...)
	}

	return all
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
func (p *perp) Signal(strategyName strategy.StrategyName) strategy.PerpSignalBuilder {
	return signal.NewPerpBuilder(strategyName, p.timeProvider)
}

// Log returns the trading logger for strategy-specific logging.
func (p *perp) Log() logging.TradingLogger {
	return p.tradingLogger
}

// PNL returns the profit and loss calculator for the perp context.
func (p *perp) PNL() perpTypes.PerpPNL {
	return p.pnl
}

var _ perpTypes.Perp = (*perp)(nil)
