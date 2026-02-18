package registry

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type pairRegistry struct {
	pairs map[string]*pairState
	mu    sync.RWMutex
}

type pairState struct {
	pair            portfolio.Pair
	instrumentTypes []connector.Instrument
}

// NewPairRegistry creates a new pair registry
func NewPairRegistry() registry.PairRegistry {
	return &pairRegistry{
		pairs: make(map[string]*pairState),
	}
}

// RegisterPair registers an pair with its supported instrument types
func (ar *pairRegistry) RegisterPair(pair portfolio.Pair, instruments ...connector.Instrument) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.pairs[pair.Symbol()] = &pairState{
		pair:            pair,
		instrumentTypes: instruments,
	}
}

// GetRequiredPairs returns all registered pairs
func (ar *pairRegistry) GetRequiredPairs() []portfolio.Pair {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	pairs := make([]portfolio.Pair, 0, len(ar.pairs))
	for _, state := range ar.pairs {
		pairs = append(pairs, state.pair)
	}

	return pairs
}

// GetPairRequirements returns all registered pairs with their instrument types
func (ar *pairRegistry) GetPairRequirements() []registry.PairRequirement {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	requirements := make([]registry.PairRequirement, 0, len(ar.pairs))
	for _, state := range ar.pairs {
		requirements = append(requirements, registry.PairRequirement{
			Pair:        state.pair,
			Instruments: state.instrumentTypes,
		})
	}

	return requirements
}

// GetInstrumentTypes returns the instrument types supported for an pair
func (ar *pairRegistry) GetInstrumentTypes(pair portfolio.Pair) []connector.Instrument {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	state, exists := ar.pairs[pair.Symbol()]
	if !exists {
		return []connector.Instrument{}
	}

	return state.instrumentTypes
}
