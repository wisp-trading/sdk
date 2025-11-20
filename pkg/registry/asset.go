package registry

import (
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

type assetRegistry struct {
	strategyRegistry registry.StrategyRegistry

	// Cached asset requirements
	requirements []registry.AssetRequirement
	assetCache   map[string]portfolio.Asset
	mu           sync.RWMutex
}

// NewAssetRegistry creates a registry that derives asset requirements from enabled strategies
func NewAssetRegistry(strategyRegistry registry.StrategyRegistry) registry.AssetRegistry {
	ar := &assetRegistry{
		strategyRegistry: strategyRegistry,
		requirements:     []registry.AssetRequirement{},
		assetCache:       make(map[string]portfolio.Asset),
	}
	// Initialize cache with current strategies
	ar.rebuildCache()
	return ar
}

func (ar *assetRegistry) GetRequiredAssets() []portfolio.Asset {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Return cached assets
	assets := make([]portfolio.Asset, 0, len(ar.assetCache))
	for _, asset := range ar.assetCache {
		assets = append(assets, asset)
	}

	return assets
}

func (ar *assetRegistry) GetAssetRequirements() []registry.AssetRequirement {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Return cached requirements
	result := make([]registry.AssetRequirement, len(ar.requirements))
	copy(result, ar.requirements)
	return result
}

func (ar *assetRegistry) GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	// Find in cached requirements
	for _, req := range ar.requirements {
		if req.Asset.Symbol() == asset.Symbol() {
			return req.Instruments
		}
	}

	return []connector.Instrument{}
}

// RefreshAssets rebuilds the asset cache from current enabled strategies
// This should be called when strategies are registered/enabled/disabled
func (ar *assetRegistry) RefreshAssets() {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.rebuildCache()
}

// rebuildCache recalculates asset requirements from enabled strategies
// Must be called with write lock held
func (ar *assetRegistry) rebuildCache() {
	enabledStrategies := ar.strategyRegistry.GetEnabledStrategies()

	// Clear caches
	ar.requirements = []registry.AssetRequirement{}
	ar.assetCache = make(map[string]portfolio.Asset)

	// Rebuild from enabled strategies
	for _, strat := range enabledStrategies {
		requiredAssets := strat.GetRequiredAssets()

		for _, req := range requiredAssets {
			ar.requirements = append(ar.requirements, registry.AssetRequirement{
				Asset:       req.Symbol,
				Instruments: req.Instruments,
			})
			ar.assetCache[req.Symbol.Symbol()] = req.Symbol
		}
	}
}
