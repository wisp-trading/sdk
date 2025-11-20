package plugin

import (
	"plugin"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
)

// LoadedPlugin represents a plugin that has been loaded into memory
type LoadedPlugin struct {
	ID           uuid.UUID
	Name         string
	Plugin       *plugin.Plugin
	StrategyFunc func() strategy.Strategy
	Metadata     *Metadata
}

// Metadata contains information about a plugin
type Metadata struct {
	ID          uuid.UUID
	Name        string
	Description string
	RiskLevel   string
	Type        string
	Version     string
	PluginPath  string
	CreatedBy   string
	Parameters  map[string]ParameterDef
}

// ParameterDef defines a strategy parameter
type ParameterDef struct {
	Name        string
	Type        string
	Description string
	Default     interface{}
	Required    bool
	Min         interface{}
	Max         interface{}
}

// ParameterProvider is an optional interface that strategies can implement
type ParameterProvider interface {
	GetParameters() []ParameterDef
}
