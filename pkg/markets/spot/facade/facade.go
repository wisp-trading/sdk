// Package facade is the domain context injected into strategies (wisp.Spot/Perp/…).
// Layout: markets/<domain>/facade — same shell for every market.
package facade

import (
	baseFacade "github.com/wisp-trading/sdk/pkg/markets/base/facade"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	spotSignal "github.com/wisp-trading/sdk/pkg/markets/spot/signal"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type spot struct {
	baseFacade.Core
	pnl spotTypes.SpotPNL
}

func NewSpot(
	tradingLogger logging.TradingLogger,
	watchlist spotTypes.SpotWatchlist,
	store spotTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	pnl spotTypes.SpotPNL,
	router execution.SignalRouter,
) spotTypes.Spot {
	return &spot{
		Core: baseFacade.Core{
			Logger:       tradingLogger,
			Watchlist:    watchlist,
			Store:        store,
			TimeProvider: timeProvider,
			Router:       router,
		},
		pnl: pnl,
	}
}

// Signal creates a new spot signal builder for the given strategy.
func (s *spot) Signal(strategyName strategy.StrategyName) spotTypes.SpotSignalBuilder {
	return spotSignal.NewSpotBuilder(strategyName, s.TimeProvider)
}

// Emit routes a spot signal to the executor (places orders).
func (s *spot) Emit(signal spotTypes.SpotSignal) execution.ExecutionCallback {
	return s.Core.Emit(signal)
}

// Positions returns placed orders (spot "positions" are open orders — shared Orders helper).
func (s *spot) Positions(q ...market.ActivityQuery) []connector.Order {
	return s.Orders(q...)
}

// PNL returns the profit and loss calculator for the spot context.
func (s *spot) PNL() spotTypes.SpotPNL {
	return s.pnl
}

var _ spotTypes.Spot = (*spot)(nil)
