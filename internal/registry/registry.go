package strategy

import (
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type StrategyRegistry interface {
	Register(strategy strategy.Strategy)
	GetStrategy(name strategy.StrategyName) (strategy.Strategy, bool)

	GetAllStrategies() []strategy.Strategy
	GetEnabledStrategies() []strategy.Strategy

	ListStrategies() []strategy.StrategyName
	ListEnabledStrategies() []strategy.StrategyName

	EnableStrategy(name strategy.StrategyName) error
	DisableStrategy(name strategy.StrategyName) error
}

type strategyRegistry struct {
	strategies map[strategy.StrategyName]strategy.Strategy
	mu         sync.RWMutex
}

func NewStrategyRegistry() StrategyRegistry {
	return &strategyRegistry{
		strategies: make(map[strategy.StrategyName]strategy.Strategy),
	}
}

func (sr *strategyRegistry) Register(strategy strategy.Strategy) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.strategies[strategy.GetName()] = strategy
}

func (sr *strategyRegistry) GetStrategy(name strategy.StrategyName) (strategy.Strategy, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	strat, exists := sr.strategies[name]
	return strat, exists
}

// GetAllStrategies returns all registered strategies (enabled and disabled)
func (sr *strategyRegistry) GetAllStrategies() []strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	strategies := make([]strategy.Strategy, 0, len(sr.strategies))
	for _, strategy := range sr.strategies {
		strategies = append(strategies, strategy)
	}
	return strategies
}

// GetEnabledStrategies returns only enabled strategies
func (sr *strategyRegistry) GetEnabledStrategies() []strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var enabled []strategy.Strategy
	for _, strategy := range sr.strategies {
		if strategy.IsEnabled() {
			enabled = append(enabled, strategy)
		}
	}
	return enabled
}

// ListStrategies returns all strategy names (kept for backward compatibility)
func (sr *strategyRegistry) ListStrategies() []strategy.StrategyName {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	names := make([]strategy.StrategyName, 0, len(sr.strategies))
	for name := range sr.strategies {
		names = append(names, name)
	}
	return names
}

// ListEnabledStrategies returns only enabled strategy names
func (sr *strategyRegistry) ListEnabledStrategies() []strategy.StrategyName {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	var names []strategy.StrategyName
	for _, strategy := range sr.strategies {
		if strategy.IsEnabled() {
			names = append(names, strategy.GetName())
		}
	}
	return names
}

// EnableStrategy enables a strategy by name
func (sr *strategyRegistry) EnableStrategy(name strategy.StrategyName) error {
	sr.mu.RLock()
	strategy, exists := sr.strategies[name]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy %s not found", name)
	}

	return strategy.Enable()
}

// DisableStrategy disables a strategy by name
func (sr *strategyRegistry) DisableStrategy(name strategy.StrategyName) error {
	sr.mu.RLock()
	strategy, exists := sr.strategies[name]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy %s not found", name)
	}

	return strategy.Disable()
}
