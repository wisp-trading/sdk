package registry

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PairRequirement defines which instrument types are needed for an asset
type PairRequirement struct {
	Asset       portfolio.Pair
	Instruments []connector.Instrument
}

// PairRegistry manages assets explicitly registered by the application
type PairRegistry interface {
	// RegisterPair registers an asset with its supported instrument types
	RegisterPair(pair portfolio.Pair, instruments ...connector.Instrument)

	// GetRequiredPairs returns all registered assets
	GetRequiredPairs() []portfolio.Pair

	// GetPairRequirements returns detailed requirements including instrument types per pair
	GetPairRequirements() []PairRequirement

	// GetInstrumentTypes returns the instrument types supported for a specific asset
	GetInstrumentTypes(asset portfolio.Pair) []connector.Instrument
}
