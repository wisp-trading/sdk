// Package facade is the domain context injected into strategies (wisp.Spot/Perp/…).
// Layout: markets/<domain>/facade — same shell for every market.
package facade

import (
	baseFacade "github.com/wisp-trading/sdk/pkg/markets/base/facade"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/perp/signal"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type perp struct {
	baseFacade.Core
	store perpTypes.MarketStore
	pnl   perpTypes.PerpPNL
}

func NewPerp(
	tradingLogger logging.TradingLogger,
	watchlist perpTypes.PerpWatchlist,
	store perpTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	pnl perpTypes.PerpPNL,
	router execution.SignalRouter,
) perpTypes.Perp {
	return &perp{
		Core: baseFacade.Core{
			Logger:       tradingLogger,
			Watchlist:    watchlist,
			Store:        store,
			TimeProvider: timeProvider,
			Router:       router,
		},
		store: store,
		pnl:   pnl,
	}
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

// Position returns a single live position for a specific exchange + pair from the store.
func (p *perp) Position(exchange connector.ExchangeName, pair portfolio.Pair) (*perpConn.Position, bool) {
	pos := p.store.GetPosition(exchange, pair)
	if pos == nil {
		return nil, false
	}
	return pos, true
}

// Positions returns live positions from the store.
func (p *perp) Positions(q ...market.ActivityQuery) []perpConn.Position {
	if len(q) > 0 {
		return p.store.QueryPositions(q[0])
	}
	return p.store.GetPositions()
}

// Signal creates a new perp signal builder for the given strategy.
func (p *perp) Signal(strategyName strategy.StrategyName) perpTypes.PerpSignalBuilder {
	return signal.NewPerpBuilder(strategyName, p.TimeProvider)
}

// Emit routes a perp signal to the executor (places orders).
func (p *perp) Emit(signal perpTypes.PerpSignal) execution.ExecutionCallback {
	return p.Core.Emit(signal)
}

// PNL returns the profit and loss calculator for the perp context.
func (p *perp) PNL() perpTypes.PerpPNL {
	return p.pnl
}

var _ perpTypes.Perp = (*perp)(nil)
