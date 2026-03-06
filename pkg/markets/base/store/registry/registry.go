package registry

import (
	"sync"

	market2 "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type marketRegistry struct {
	mu     sync.RWMutex
	stores map[connector.MarketType]market2.MarketStore
}

func NewMarketRegistry() market2.MarketRegistry {
	return &marketRegistry{
		stores: make(map[connector.MarketType]market2.MarketStore),
	}
}

func (r *marketRegistry) Get(marketType connector.MarketType) market2.MarketStore {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.stores[marketType]
}

func (r *marketRegistry) GetAll() map[connector.MarketType]market2.MarketStore {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[connector.MarketType]market2.MarketStore, len(r.stores))
	for k, v := range r.stores {
		result[k] = v
	}
	return result
}

func (r *marketRegistry) Register(store market2.MarketStore) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stores[store.MarketType()] = store
}

func (r *marketRegistry) Types() []connector.MarketType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]connector.MarketType, 0, len(r.stores))
	for t := range r.stores {
		types = append(types, t)
	}
	return types
}
