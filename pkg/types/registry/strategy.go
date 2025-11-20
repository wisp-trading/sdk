package registry

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// StrategyRegistry manages strategy registration and lifecycle
type StrategyRegistry interface {
	// GetStrategy retrieves a strategy by name
	GetStrategy(name strategy.StrategyName) (strategy.Strategy, bool)

	// RegisterStrategy registers a single strategy
	RegisterStrategy(strat strategy.Strategy)

	// RegisterAllStrategies registers multiple strategies at once
	RegisterAllStrategies(strategies []strategy.Strategy)

	// GetAllStrategies returns all registered strategies
	GetAllStrategies() []strategy.Strategy

	// GetEnabledStrategies returns only enabled strategies
	GetEnabledStrategies() []strategy.Strategy

	// EnableStrategy enables a strategy by name
	EnableStrategy(name strategy.StrategyName) error

	// DisableStrategy disables a strategy by name
	DisableStrategy(name strategy.StrategyName) error

	// IsStrategyEnabled checks if a strategy is enabled
	IsStrategyEnabled(name strategy.StrategyName) bool

	// GetStrategyCount returns the total number of registered strategies
	GetStrategyCount() int
}
