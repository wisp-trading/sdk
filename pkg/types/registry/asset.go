package registry

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// AssetRequirement defines which instrument types are needed for an asset
type AssetRequirement struct {
	Asset       portfolio.Asset
	Instruments []connector.Instrument
}

// AssetRegistry manages assets explicitly registered by the application
type AssetRegistry interface {
	// RegisterAsset registers an asset with its supported instrument types
	RegisterAsset(asset portfolio.Asset, instruments ...connector.Instrument)

	// GetRequiredAssets returns all registered assets
	GetRequiredAssets() []portfolio.Asset

	// GetAssetRequirements returns detailed requirements including instrument types per asset
	GetAssetRequirements() []AssetRequirement

	// GetInstrumentTypes returns the instrument types supported for a specific asset
	GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument
}
