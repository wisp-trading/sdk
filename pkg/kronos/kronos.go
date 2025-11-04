package kronos

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
)

// Kronos is the base context object for strategy GetSignals methods.
// It provides read-only access to market data, indicators, and analytics.
// It does NOT include trade execution methods.
type Kronos struct {
	store  store.Store
	logger logging.ApplicationLogger

	// Namespaced services for user-friendly API
	Indicators *IndicatorService
	Market     *MarketService
	Analytics  *AnalyticsService
}

// NewKronos creates a new Kronos context with the given store and logger.
// This is injected via fx DI into strategies.
func NewKronos(store store.Store, logger logging.ApplicationLogger) *Kronos {
	k := &Kronos{
		store:  store,
		logger: logger,
	}

	// Initialize services with references to store and logger
	k.Indicators = &IndicatorService{
		store:  store,
		logger: logger,
	}
	k.Market = &MarketService{
		store:  store,
		logger: logger,
	}
	k.Analytics = &AnalyticsService{
		store:  store,
		logger: logger,
	}

	return k
}

// Log returns the logger for strategy logging.
// Usage: k.Log().Info("checking signals")
func (k *Kronos) Log() logging.ApplicationLogger {
	return k.logger
}

// Store returns the underlying store for advanced use cases.
// Most users should use the service methods instead.
func (k *Kronos) Store() store.Store {
	return k.store
}

// Asset creates a new portfolio.Asset from a symbol string.
// This is a convenience method to avoid importing portfolio everywhere.
// Usage: btc := k.Asset("BTC")
func (k *Kronos) Asset(symbol string) portfolio.Asset {
	return portfolio.NewAsset(symbol)
}
