package registry

import (
	"sync"

	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
	"github.com/wisp-trading/wisp/pkg/types/registry"
)

type assetRegistry struct {
	assets map[string]*assetState
	mu     sync.RWMutex
}

type assetState struct {
	asset           portfolio.Asset
	instrumentTypes []connector.Instrument
}

// NewAssetRegistry creates a new asset registry
func NewAssetRegistry() registry.AssetRegistry {
	return &assetRegistry{
		assets: make(map[string]*assetState),
	}
}

// RegisterAsset registers an asset with its supported instrument types
func (ar *assetRegistry) RegisterAsset(asset portfolio.Asset, instruments ...connector.Instrument) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.assets[asset.Symbol()] = &assetState{
		asset:           asset,
		instrumentTypes: instruments,
	}
}

// GetRequiredAssets returns all registered assets
func (ar *assetRegistry) GetRequiredAssets() []portfolio.Asset {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	assets := make([]portfolio.Asset, 0, len(ar.assets))
	for _, state := range ar.assets {
		assets = append(assets, state.asset)
	}

	return assets
}

// GetAssetRequirements returns all registered assets with their instrument types
func (ar *assetRegistry) GetAssetRequirements() []registry.AssetRequirement {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	requirements := make([]registry.AssetRequirement, 0, len(ar.assets))
	for _, state := range ar.assets {
		requirements = append(requirements, registry.AssetRequirement{
			Asset:       state.asset,
			Instruments: state.instrumentTypes,
		})
	}

	return requirements
}

// GetInstrumentTypes returns the instrument types supported for an asset
func (ar *assetRegistry) GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	state, exists := ar.assets[asset.Symbol()]
	if !exists {
		return []connector.Instrument{}
	}

	return state.instrumentTypes
}
