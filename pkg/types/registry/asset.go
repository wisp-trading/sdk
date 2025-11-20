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

// AssetRegistry manages required assets based on enabled strategies
type AssetRegistry interface {
	// GetRequiredAssets returns all assets needed by enabled strategies
	GetRequiredAssets() []portfolio.Asset

	// GetAssetRequirements returns detailed requirements including instrument types per asset
	GetAssetRequirements() []AssetRequirement

	// GetInstrumentTypes returns the instrument types needed for a specific asset
	GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument

	// RefreshAssets rebuilds the asset cache from current enabled strategies
	// Should be called when strategies are registered/enabled/disabled
	RefreshAssets()
}
