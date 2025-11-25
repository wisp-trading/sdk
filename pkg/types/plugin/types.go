package plugin

import (
	"plugin"

	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
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

// LoadedHookPlugin represents a hook plugin that has been loaded into memory
type LoadedHookPlugin struct {
	ID         uuid.UUID
	Name       string
	Plugin     *plugin.Plugin
	HookPlugin execution.HookPlugin
	Metadata   *Metadata
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
	// StrategyPlugin indicates a strategy plugin
	StrategyPlugin PluginType = iota

	// HookPlugin indicates a hook plugin
	HookPlugin
)
