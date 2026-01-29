package plugin

import (
	"plugin"

	"github.com/google/uuid"
	"github.com/wisp-trading/wisp/pkg/types/strategy"
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
	SDKVersion  string     // SDK version the plugin was built against
	PluginType  PluginType // Type of plugin (strategy or hook)
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

// PluginType indicates the type of plugin (strategy or hook)
type PluginType int

const (
	// UnknownPlugin indicates an uninitialized or invalid plugin type
	UnknownPlugin PluginType = iota

	// StrategyPlugin indicates a strategy plugin
	StrategyPlugin

	// HookPlugin indicates a hook plugin
	HookPlugin
)
