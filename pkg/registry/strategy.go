package registry

import (
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type strategyRegistry struct {
	strategies map[strategy.StrategyName]strategy.Strategy
	mu         sync.RWMutex
}

// NewStrategyRegistry creates a new strategy registry
func NewStrategyRegistry() registry.StrategyRegistry {
	return &strategyRegistry{
		strategies: make(map[strategy.StrategyName]strategy.Strategy),
	}
}

func (sr *strategyRegistry) GetStrategy(name strategy.StrategyName) (strategy.Strategy, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	strat, exists := sr.strategies[name]
	return strat, exists
}

func (sr *strategyRegistry) RegisterStrategy(strat strategy.Strategy) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.strategies[strat.GetName()] = strat
}

func (sr *strategyRegistry) RegisterAllStrategies(strategies []strategy.Strategy) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, strat := range strategies {
		sr.strategies[strat.GetName()] = strat
	}
}

func (sr *strategyRegistry) GetAllStrategies() []strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	strategies := make([]strategy.Strategy, 0, len(sr.strategies))
	for _, strat := range sr.strategies {
		strategies = append(strategies, strat)
	}

	return strategies
}

func (sr *strategyRegistry) GetEnabledStrategies() []strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	enabled := make([]strategy.Strategy, 0)
	for _, strat := range sr.strategies {
		if strat.IsEnabled() {
			enabled = append(enabled, strat)
		}
	}

	return enabled
}

func (sr *strategyRegistry) EnableStrategy(name strategy.StrategyName) error {
	sr.mu.RLock()
	strat, exists := sr.strategies[name]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy %s not found", name)
	}

	return strat.Enable()
}

func (sr *strategyRegistry) DisableStrategy(name strategy.StrategyName) error {
	sr.mu.RLock()
	strat, exists := sr.strategies[name]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy %s not found", name)
	}

	return strat.Disable()
}

func (sr *strategyRegistry) IsStrategyEnabled(name strategy.StrategyName) bool {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	strat, exists := sr.strategies[name]
	if !exists {
		return false
	}

	return strat.IsEnabled()
}

func (sr *strategyRegistry) GetStrategyCount() int {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return len(sr.strategies)
}
