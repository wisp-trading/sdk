package wisp

import (
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
	"github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// wisp is the SDK context object injected into strategies.
// Order placement is only via market domains: Spot/Perp/Predict/Options/Onchain.Emit.
type wisp struct {
	tradingLogger logging.TradingLogger
	indicators    analytics.Indicators
	analytics     analytics.Analytics
	activity      activity.Activity
	perp          perpTypes.Perp
	predict       predTypes.Predict
	spotService   spotTypes.Spot
	options       optionsTypes.Options
	onchain       onchainTypes.Onchain
	priceFeeds    types.PriceFeeds
}

// NewWisp creates a new Wisp context with injected services.
func NewWisp(
	tradingLogger logging.TradingLogger,
	indicators analytics.Indicators,
	analyticsService analytics.Analytics,
	activityService activity.Activity,
	perpService perpTypes.Perp,
	predictService predTypes.Predict,
	spotService spotTypes.Spot,
	optionsService optionsTypes.Options,
	onchainService onchainTypes.Onchain,
	priceFeeds types.PriceFeeds,
) wispTypes.Wisp {
	return &wisp{
		tradingLogger: tradingLogger,
		indicators:    indicators,
		analytics:     analyticsService,
		activity:      activityService,
		perp:          perpService,
		predict:       predictService,
		spotService:   spotService,
		options:       optionsService,
		onchain:       onchainService,
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
func (k *wisp) Onchain() onchainTypes.Onchain    { return k.onchain }
func (k *wisp) PriceFeeds() types.PriceFeeds     { return k.priceFeeds }

func (k *wisp) Pair(base, quote portfolio.Asset) portfolio.Pair {
	return portfolio.NewPair(base, quote)
}

func (k *wisp) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}
