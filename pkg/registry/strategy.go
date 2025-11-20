package registry

import (
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// StrategyRegistry manages loaded strategy instances
type StrategyRegistry interface {
	// Register adds a strategy instance to the registry
	// id: unique identifier for this instance
	Register(id string, strat strategy.Strategy) error

	// Unregister removes a strategy instance
	Unregister(id string) error

	// Get retrieves a strategy instance by ID
	Get(id string) (strategy.Strategy, error)

	// GetByName retrieves the first strategy with the given name
	GetByName(name strategy.StrategyName) (strategy.Strategy, error)

	// List returns all registered strategy IDs
	List() []string

	// ListByName returns IDs of strategies with the given name
	ListByName(name strategy.StrategyName) []string

	// GetAll returns all registered strategies
	GetAll() map[string]strategy.Strategy

	// GetEnabled returns only enabled strategy instances
	GetEnabled() map[string]strategy.Strategy

	// Enable enables a strategy by ID
	Enable(id string) error

	// Disable disables a strategy by ID
	Disable(id string) error

	// IsEnabled checks if a strategy is enabled
	IsEnabled(id string) bool

	// Count returns the number of registered strategies
	Count() int
}

type strategyRegistry struct {
	strategies map[string]strategy.Strategy
	mu         sync.RWMutex
}

// NewStrategyRegistry creates a new strategy registry
func NewStrategyRegistry() StrategyRegistry {
	return &strategyRegistry{
		strategies: make(map[string]strategy.Strategy),
	}
}

func (sr *strategyRegistry) Register(id string, strat strategy.Strategy) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if _, exists := sr.strategies[id]; exists {
		return fmt.Errorf("strategy with ID %s already registered", id)
	}

	sr.strategies[id] = strat
	return nil
}

func (sr *strategyRegistry) Unregister(id string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if _, exists := sr.strategies[id]; !exists {
		return fmt.Errorf("strategy with ID %s not found", id)
	}

	delete(sr.strategies, id)
	return nil
}

func (sr *strategyRegistry) Get(id string) (strategy.Strategy, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	strat, exists := sr.strategies[id]
	if !exists {
		return nil, fmt.Errorf("strategy with ID %s not found", id)
	}

	return strat, nil
}

func (sr *strategyRegistry) GetByName(name strategy.StrategyName) (strategy.Strategy, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	for _, strat := range sr.strategies {
		if strat.GetName() == name {
			return strat, nil
		}
	}

	return nil, fmt.Errorf("strategy with name %s not found", name)
}

func (sr *strategyRegistry) List() []string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	ids := make([]string, 0, len(sr.strategies))
	for id := range sr.strategies {
		ids = append(ids, id)
	}

	return ids
}

func (sr *strategyRegistry) ListByName(name strategy.StrategyName) []string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	ids := make([]string, 0)
	for id, strat := range sr.strategies {
		if strat.GetName() == name {
			ids = append(ids, id)
		}
	}

	return ids
}

func (sr *strategyRegistry) GetAll() map[string]strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	// Return a copy to prevent external modification
	strategies := make(map[string]strategy.Strategy, len(sr.strategies))
	for id, strat := range sr.strategies {
		strategies[id] = strat
	}

	return strategies
}

func (sr *strategyRegistry) GetEnabled() map[string]strategy.Strategy {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	enabled := make(map[string]strategy.Strategy)
	for id, strat := range sr.strategies {
		if strat.IsEnabled() {
			enabled[id] = strat
		}
	}

	return enabled
}

func (sr *strategyRegistry) Enable(id string) error {
	sr.mu.RLock()
	strat, exists := sr.strategies[id]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy with ID %s not found", id)
	}

	return strat.Enable()
}

func (sr *strategyRegistry) Disable(id string) error {
	sr.mu.RLock()
	strat, exists := sr.strategies[id]
	sr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("strategy with ID %s not found", id)
	}

	return strat.Disable()
}

func (sr *strategyRegistry) IsEnabled(id string) bool {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	if strat, exists := sr.strategies[id]; exists {
		return strat.IsEnabled()
	}

	return false
}

func (sr *strategyRegistry) Count() int {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return len(sr.strategies)
}
