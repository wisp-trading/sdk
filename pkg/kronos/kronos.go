package kronos

import (
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/indicators"
	marketService "github.com/backtesting-org/kronos-sdk/pkg/kronos/market"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/signal"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// Kronos is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
type Kronos struct {
	store         market.MarketData
	tradingLogger logging.TradingLogger

	// Namespaced services for user-friendly API
	Indicators *indicators.IndicatorService
	Market     *marketService.MarketService
	Analytics  *analytics.AnalyticsService
	signalSvc  *signal.Service
}

// NewKronos creates a new Kronos context with injected services.
// This is injected via fx DI into strategies.
func NewKronos(
	store market.MarketData,
	tradingLogger logging.TradingLogger,
	indicators *indicators.IndicatorService,
	market *marketService.MarketService,
	analytics *analytics.AnalyticsService,
	signalSvc *signal.Service,
) *Kronos {
	return &Kronos{
		store:         store,
		tradingLogger: tradingLogger,
		Indicators:    indicators,
		Market:        market,
		Analytics:     analytics,
		signalSvc:     signalSvc,
	}
}

// Log returns the trading logger for strategy logging.
// Usage: k.Log().Opportunity("MyStrategy", "BTC", "Found signal")
func (k *Kronos) Log() logging.TradingLogger {
	return k.tradingLogger
}

// Store returns the underlying store for advanced use cases.
// Most users should use the service methods instead.
func (k *Kronos) Store() market.MarketData {
	return k.store
}

// Asset creates a new portfolio.Asset from a symbol string.
// This is a convenience method to avoid importing portfolio everywhere.
// Usage: btc := k.Asset("BTC")
func (k *Kronos) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}

// Signal creates a new signal builder for constructing trading signals.
// Usage: k.Signal(strategyName).Buy(asset, exchange, qty).SellShort(asset, exchange, qty).Build()
func (k *Kronos) Signal(strategyName strategy.StrategyName) *signal.Builder {
	return k.signalSvc.New(strategyName)
}
