package wisp

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
	"github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// wisp is the SDK context object injected into strategies.
// It provides access to market data, indicators, analytics, and signal dispatch.
type wisp struct {
	tradingLogger    logging.TradingLogger
	universeProvider UniverseProvider
	indicators       analytics.Indicators
	analytics        analytics.Analytics
	market           analytics.Market
	signal           strategy.SignalFactory
	activity         activity.Activity
	router           execution.SignalRouter
	perp             perpTypes.Perp
	predict          predTypes.Predict
}

// NewWisp creates a new Wisp context with injected services.
// This is injected via fx DI into strategies.
func NewWisp(
	tradingLogger logging.TradingLogger,
	universeProvider UniverseProvider,
	indicators analytics.Indicators,
	analyticsService analytics.Analytics,
	marketService analytics.Market,
	signal strategy.SignalFactory,
	activityService activity.Activity,
	router execution.SignalRouter,
	perpService perpTypes.Perp,
	predictService predTypes.Predict,
) wispTypes.Wisp {
	return &wisp{
		tradingLogger:    tradingLogger,
		universeProvider: universeProvider,
		indicators:       indicators,
		market:           marketService,
		analytics:        analyticsService,
		signal:           signal,
		activity:         activityService,
		router:           router,
		perp:             perpService,
		predict:          predictService,
	}
}

func (k *wisp) Indicators() analytics.Indicators { return k.indicators }
func (k *wisp) Analytics() analytics.Analytics   { return k.analytics }
func (k *wisp) Market() analytics.Market         { return k.market }
func (k *wisp) Activity() activity.Activity      { return k.activity }
func (k *wisp) Log() logging.TradingLogger       { return k.tradingLogger }
func (k *wisp) Perp() perpTypes.Perp             { return k.perp }
func (k *wisp) Predict() predTypes.Predict       { return k.predict }

func (k *wisp) Pair(base, quote portfolio.Asset) portfolio.Pair {
	return portfolio.NewPair(base, quote)
}

func (k *wisp) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}

// SpotSignal creates a new signal builder for spot market trading signals.
func (k *wisp) SpotSignal(strategyName strategy.StrategyName) strategy.SpotSignalBuilder {
	return k.signal.NewSpot(strategyName)
}

// Universe returns the live spot trading universe.
// Provides the complete trading universe available to the strategy.
// Data is cached since it does not change after initialization.
func (k *wisp) Universe() wispTypes.Universe {
	return k.universeProvider.Universe()
}

// Emit routes a signal directly to the executor via the SDK's SignalRouter.
// This is the primary way strategies dispatch signals — non-blocking, fire and forget.
func (k *wisp) Emit(signal strategy.Signal) {
	k.router.Route(signal)
}
