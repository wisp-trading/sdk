package ingestors

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// AssetRequirement defines which instrument types are needed for an asset
type AssetRequirement struct {
	Asset       portfolio.Asset
	Instruments []connector.Instrument // e.g., [TypeSpot, TypePerpetual]
}

type AssetInterest interface {
	GetRequiredAssets() []portfolio.Asset
	IsAssetRequired(symbol string) bool

	// Get detailed requirements including instrument types
	GetAssetRequirements() []AssetRequirement
	GetInstrumentTypes(asset portfolio.Asset) []connector.Instrument
}
