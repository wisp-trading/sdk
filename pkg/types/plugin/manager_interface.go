package plugin

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// Manager defines the interface for plugin management operations
type Manager interface {
	// LoadStrategyPlugin loads a strategy plugin and registers it with the registry
	// Returns the loaded strategy instance for reference (though it's already registered)
	LoadStrategyPlugin(pluginPath string) (strategy.Strategy, error)

	// LoadHookPlugin loads a hook plugin and registers hooks with the registry
	LoadHookPlugin(pluginPath string) error
}
