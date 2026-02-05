package wisp

import (
	"github.com/wisp-trading/sdk/pkg/inference/features"
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
	tradingLogger     logging.TradingLogger
	universeProvider  UniverseProvider
	indicators        analytics.Indicators
	analytics         analytics.Analytics
	market            analytics.Market
	signal            strategy.SignalFactory
	featureAggregator features.FeatureAggregator
	activity          activity.Activity
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
	featureAggregator features.FeatureAggregator,
	activityService activity.Activity,
) wispTypes.Wisp {
	return &wisp{
		tradingLogger:     tradingLogger,
		universeProvider:  universeProvider,
		indicators:        indicators,
		market:            marketService,
		analytics:         analyticsService,
		signal:            signal,
		featureAggregator: featureAggregator,
		activity:          activityService,
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

// Asset creates a new portfolio.Pair from a symbol string.
// This is a convenience method to avoid importing portfolio everywhere.
// Usage: btc := k.Pair("BTC")
func (k *wisp) Asset(symbol string) portfolio.Pair {
	return portfolio.NewAsset(symbol)
}

// Signal creates a new signal builder for constructing trading signals.
// Usage: k.Signal(strategyName).Buy(asset, exchange, qty).SellShort(asset, exchange, qty).Build()
func (k *wisp) Signal(strategyName strategy.StrategyName) strategy.SignalBuilder {
	return k.signal.New(strategyName)
}

// Universe returns the tradeable assets, instruments, and ready exchanges.
// Provides the complete trading universe available to the strategy.
// Data is cached since it does not change after initialization.
func (k *wisp) Universe() wispTypes.Universe {
	return k.universeProvider.Universe()
}

// Features returns the ML feature aggregator for extracting market features.
// Provides access to 41+ features including market data, orderbook, technical indicators,
// volatility, volume, price metrics, and time-based features.
// Usage: featureMap, err := k.Features().Extract(asset)
func (k *wisp) Features() features.FeatureAggregator {
	return k.featureAggregator
}
