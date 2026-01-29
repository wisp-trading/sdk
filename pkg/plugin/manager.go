package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/wisp-trading/wisp/pkg/types/execution"
	"github.com/wisp-trading/wisp/pkg/types/logging"
	plugintypes "github.com/wisp-trading/wisp/pkg/types/plugin"
	"github.com/wisp-trading/wisp/pkg/types/registry"
	"github.com/wisp-trading/wisp/pkg/types/strategy"
	"github.com/wisp-trading/wisp/pkg/types/wisp"
)

// manager is the unexported implementation of plugintypes.Manager
type manager struct {
	logger           logging.ApplicationLogger
	hookRegistry     registry.Hooks
	strategyRegistry registry.StrategyRegistry
	wisp             wisp.Wisp
}

// NewManager creates a new plugin manager
func NewManager(
	logging logging.ApplicationLogger,
	hookRegistry registry.Hooks,
	strategyRegistry registry.StrategyRegistry,
	wisp wisp.Wisp,
) plugintypes.Manager {
	return &manager{
		logger:           logging,
		hookRegistry:     hookRegistry,
		strategyRegistry: strategyRegistry,
		wisp:             wisp,
	}
}

// LoadStrategyPlugin loads a strategy plugin and registers it with the strategy registry
func (m *manager) LoadStrategyPlugin(pluginPath string) (strategy.Strategy, error) {
	m.logger.Info("Loading strategy plugin", "path", pluginPath)

	// Validate file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin file does not exist: %s", pluginPath)
	}

	// Extract and validate SDK version
	pluginSDKVersion, err := extractSDKVersionFromPath(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract SDK version: %w", err)
	}

	if err := validateSDKVersion(pluginSDKVersion); err != nil {
		return nil, err
	}

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the NewStrategy symbol
	newStrategySymbol, err := p.Lookup("NewStrategy")
	if err != nil {
		return nil, fmt.Errorf("plugin must export NewStrategy function: %w", err)
	}

	// Type assert to constructor function
	newStrategyFunc, ok := newStrategySymbol.(func(wisp.Wisp) strategy.Strategy)
	if !ok {
		return nil, fmt.Errorf("NewStrategy must have signature: func(wisp.Wisp) strategy.Strategy, got %T", newStrategySymbol)
	}

	// Call it to get the strategy instance
	strat := newStrategyFunc(m.wisp)
	if strat == nil {
		return nil, fmt.Errorf("NewStrategy() returned nil")
	}

	// Register with strategy registry
	m.strategyRegistry.RegisterStrategy(strat)

	// The strategy name comes from strat.GetName(), not from the symbol
	m.logger.Info("Strategy plugin loaded and registered", "name", strat.GetName())

	return strat, nil
}

// LoadHookPlugin loads a hook plugin and registers hooks with the hook registry
func (m *manager) LoadHookPlugin(pluginPath string) error {
	m.logger.Info("Loading hook plugin", "path", pluginPath)

	// Validate file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("hook plugin file does not exist: %s", pluginPath)
	}

	// Extract and validate SDK version
	pluginSDKVersion, err := extractSDKVersionFromPath(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to extract SDK version: %w", err)
	}

	if err := validateSDKVersion(pluginSDKVersion); err != nil {
		return err
	}

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open hook plugin: %w", err)
	}

	// Look up the HookPlugin symbol
	hookPluginSymbol, err := p.Lookup("HookPlugin")
	if err != nil {
		return fmt.Errorf("hook plugin must export HookPlugin symbol: %w", err)
	}

	// Type assert to execution.HookPlugin
	hookPlugin, ok := hookPluginSymbol.(execution.HookPlugin)
	if !ok {
		return fmt.Errorf("HookPlugin symbol is not of type execution.HookPlugin")
	}

	// Instantiate and register hooks immediately
	hooks := hookPlugin.CreateHooks()
	if hooks != nil && len(hooks) > 0 {
		m.hookRegistry.RegisterHooks(hooks)
		m.logger.Info("Registered hooks from plugin", "count", len(hooks), "name", hookPlugin.Name())
	}

	return nil
}
