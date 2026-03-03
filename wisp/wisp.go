package wisp

import (
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
	"github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// Wisp is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
type wisp struct {
	tradingLogger    logging.TradingLogger
	universeProvider UniverseProvider
	indicators       analytics.Indicators
	analytics        analytics.Analytics
	market           analytics.Market
	signal           strategy.SignalFactory
	activity         activity.Activity
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
) wispTypes.Wisp {
	return &wisp{
		tradingLogger:    tradingLogger,
		universeProvider: universeProvider,
		indicators:       indicators,
		market:           marketService,
		analytics:        analyticsService,
		signal:           signal,
		activity:         activityService,
	}
}

func (k *wisp) Indicators() analytics.Indicators {
	return k.indicators
}

func (k *wisp) Analytics() analytics.Analytics {
	return k.analytics
}

func (k *wisp) Market() analytics.Market {
	return k.market
}

func (k *wisp) Activity() activity.Activity {
	return k.activity
}

// Log returns the trading logger for strategy logging.
// Usage: k.Log().Opportunity("MyStrategy", "BTC", "Found signal")
func (k *wisp) Log() logging.TradingLogger {
	return k.tradingLogger
}

// Pair creates a new portfolio.Pair from a symbol string.
// Convenience method to avoid importing portfolio package in strategies.
// Example: btc := k.Pair(base, quote)
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

// PerpSignal creates a new signal builder for perpetual futures trading signals.
func (k *wisp) PerpSignal(strategyName strategy.StrategyName) strategy.PerpSignalBuilder {
	return k.signal.NewPerp(strategyName)
}

// Universe returns the tradeable assets, instruments, and ready exchanges.
// Provides the complete trading universe available to the strategy.
// Data is cached since it does not change after initialization.
func (k *wisp) Universe() wispTypes.Universe {
	return k.universeProvider.Universe()
}
