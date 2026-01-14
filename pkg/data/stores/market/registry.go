package market

import (
	"sync"

	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
)

// marketRegistry is the concrete implementation of MarketRegistry
type marketRegistry struct {
	mu     sync.RWMutex
	stores map[marketTypes.MarketType]marketTypes.MarketStore
}

// NewMarketRegistry creates a new market registry
func NewMarketRegistry() marketTypes.MarketRegistry {
	return &marketRegistry{
		stores: make(map[marketTypes.MarketType]marketTypes.MarketStore),
	}
}

// Get returns the store for a specific market type
func (r *marketRegistry) Get(marketType marketTypes.MarketType) marketTypes.MarketStore {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.stores[marketType]
}

// GetAll returns all registered market stores
func (r *marketRegistry) GetAll() map[marketTypes.MarketType]marketTypes.MarketStore {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[marketTypes.MarketType]marketTypes.MarketStore, len(r.stores))
	for k, v := range r.stores {
		result[k] = v
	}
	return result
}

// Register adds a market store to the registry
func (r *marketRegistry) Register(store marketTypes.MarketStore) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stores[store.MarketType()] = store
}

// Types returns all registered market types
func (r *marketRegistry) Types() []marketTypes.MarketType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]marketTypes.MarketType, 0, len(r.stores))
	for t := range r.stores {
		types = append(types, t)
	}
	return types
}
