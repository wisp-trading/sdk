package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
)

// OptionsAssetLoader discovers available options expirations for underlyings
type OptionsAssetLoader interface {
	baseTypes.AssetLoader // Embeds Load(cfg *config.StartupConfig) error
}

// UniverseProvider builds the current trading universe
type OptionsUniverseProvider interface {
	Universe() OptionsUniverse
}
