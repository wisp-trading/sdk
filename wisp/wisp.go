package wisp

import (
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types"
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
	tradingLogger logging.TradingLogger
	indicators    analytics.Indicators
	analytics     analytics.Analytics
	activity      activity.Activity
	router        execution.SignalRouter
	perp          perpTypes.Perp
	predict       predTypes.Predict
	spotService   spotTypes.Spot
	options       optionsTypes.Options
	priceFeeds    types.PriceFeeds
}

// NewWisp creates a new Wisp context with injected services.
// This is injected via fx DI into strategies.
func NewWisp(
	tradingLogger logging.TradingLogger,
	indicators analytics.Indicators,
	analyticsService analytics.Analytics,
	activityService activity.Activity,
	router execution.SignalRouter,
	perpService perpTypes.Perp,
	predictService predTypes.Predict,
	spotService spotTypes.Spot,
	optionsService optionsTypes.Options,
	priceFeeds types.PriceFeeds,
) wispTypes.Wisp {
	return &wisp{
		tradingLogger: tradingLogger,
		indicators:    indicators,
		analytics:     analyticsService,
		activity:      activityService,
		router:        router,
		perp:          perpService,
		predict:       predictService,
		spotService:   spotService,
		options:       optionsService,
		priceFeeds:    priceFeeds,
	}
}

func (k *wisp) Indicators() analytics.Indicators { return k.indicators }
func (k *wisp) Analytics() analytics.Analytics   { return k.analytics }
func (k *wisp) Activity() activity.Activity      { return k.activity }
func (k *wisp) Log() logging.TradingLogger       { return k.tradingLogger }
func (k *wisp) Perp() perpTypes.Perp             { return k.perp }
func (k *wisp) Predict() predTypes.Predict       { return k.predict }
func (k *wisp) Spot() spotTypes.Spot             { return k.spotService }
func (k *wisp) Options() optionsTypes.Options    { return k.options }
func (k *wisp) PriceFeeds() types.PriceFeeds     { return k.priceFeeds }

func (k *wisp) Pair(base, quote portfolio.Asset) portfolio.Pair {
	return portfolio.NewPair(base, quote)
}

func (k *wisp) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}

// Emit routes a signal directly to the executor via the SDK's SignalRouter.
// This is the primary way strategies dispatch signals — non-blocking, fire and forget.
func (k *wisp) Emit(signal strategy.Signal) {
	k.router.Route(signal)
}
