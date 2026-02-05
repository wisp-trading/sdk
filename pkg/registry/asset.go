package registry

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type assetRegistry struct {
	assets map[string]*assetState
	mu     sync.RWMutex
}

type assetState struct {
	asset           portfolio.Pair
	instrumentTypes []connector.Instrument
}

// NewAssetRegistry creates a new asset registry
func NewAssetRegistry() registry.PairRegistry {
	return &assetRegistry{
		assets: make(map[string]*assetState),
	}
}

// RegisterAsset registers an asset with its supported instrument types
func (ar *assetRegistry) RegisterPair(asset portfolio.Pair, instruments ...connector.Instrument) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.assets[asset.Symbol()] = &assetState{
		asset:           asset,
		instrumentTypes: instruments,
	}
}

// GetRequiredAssets returns all registered assets
func (ar *assetRegistry) GetRequiredPairs() []portfolio.Pair {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	assets := make([]portfolio.Pair, 0, len(ar.assets))
	for _, state := range ar.assets {
		assets = append(assets, state.asset)
	}

	return assets
}

// GetAssetRequirements returns all registered assets with their instrument types
func (ar *assetRegistry) GetPairRequirements() []registry.PairRequirement {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	requirements := make([]registry.PairRequirement, 0, len(ar.assets))
	for _, state := range ar.assets {
		requirements = append(requirements, registry.PairRequirement{
			Asset:       state.asset,
			Instruments: state.instrumentTypes,
		})
	}

	return requirements
}

// GetInstrumentTypes returns the instrument types supported for an asset
func (ar *assetRegistry) GetInstrumentTypes(asset portfolio.Pair) []connector.Instrument {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	state, exists := ar.assets[asset.Symbol()]
	if !exists {
		return []connector.Instrument{}
	}

	return state.instrumentTypes
}
