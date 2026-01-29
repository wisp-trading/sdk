package registry

import (
	"github.com/wisp-trading/wisp/pkg/types/strategy"
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

	// GetStrategyCount returns the total number of registered strategies
	GetStrategyCount() int
}
