// Package facade is the domain context injected into strategies (wisp.Onchain).
package facade

import (
	baseFacade "github.com/wisp-trading/sdk/pkg/markets/base/facade"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	onchainSignal "github.com/wisp-trading/sdk/pkg/markets/onchain/signal"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type onchain struct {
	baseFacade.Core
	pnl onchainTypes.OnchainPNL
}

func NewOnchain(
	tradingLogger logging.TradingLogger,
	watchlist onchainTypes.OnchainWatchlist,
	store onchainTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	pnl onchainTypes.OnchainPNL,
	router execution.SignalRouter,
) onchainTypes.Onchain {
	return &onchain{
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

func (o *onchain) Signal(strategyName strategy.StrategyName) onchainTypes.OnchainSignalBuilder {
	return onchainSignal.NewOnchainBuilder(strategyName, o.TimeProvider)
}

func (o *onchain) Emit(signal onchainTypes.OnchainSignal) execution.ExecutionCallback {
	return o.Core.Emit(signal)
}

func (o *onchain) Positions(q ...market.ActivityQuery) []connector.Order {
	return o.Orders(q...)
}

func (o *onchain) PNL() onchainTypes.OnchainPNL {
	return o.pnl
}

var _ onchainTypes.Onchain = (*onchain)(nil)
