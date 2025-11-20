package kronos

import (
	marketService "github.com/backtesting-org/kronos-sdk/kronos/market"
	kronosTypes "github.com/backtesting-org/kronos-sdk/pkg/types/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// Kronos is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
type kronos struct {
	store         market.MarketData
	tradingLogger logging.TradingLogger

	// Namespaced services for user-friendly API
	Indicators analytics.Indicators
	Analytics  analytics.Analytics

	Market *marketService.MarketService
	signal strategy.SignalFactory
}

// NewKronos creates a new Kronos context with injected services.
// This is injected via fx DI into strategies.
func NewKronos(
	store market.MarketData,
	tradingLogger logging.TradingLogger,
	indicators analytics.Indicators,
	analytics analytics.Analytics,
	market *marketService.MarketService,
	signal strategy.SignalFactory,
) kronosTypes.Kronos {
	return &kronos{
		store:         store,
		tradingLogger: tradingLogger,
		Indicators:    indicators,
		Market:        market,
		Analytics:     analytics,
		signal:        signal,
	}
}

// Log returns the trading logger for strategy logging.
// Usage: k.Log().Opportunity("MyStrategy", "BTC", "Found signal")
func (k *kronos) Log() logging.TradingLogger {
	return k.tradingLogger
}

// Store returns the underlying store for advanced use cases.
// Most users should use the service methods instead.
func (k *kronos) Store() market.MarketData {
	return k.store
}

// Asset creates a new portfolio.Asset from a symbol string.
// This is a convenience method to avoid importing portfolio everywhere.
// Usage: btc := k.Asset("BTC")
func (k *kronos) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}

// Signal creates a new signal builder for constructing trading signals.
// Usage: k.Signal(strategyName).Buy(asset, exchange, qty).SellShort(asset, exchange, qty).Build()
func (k *kronos) Signal(strategyName strategy.StrategyName) strategy.SignalBuilder {
	return k.signal.New(strategyName)
}
