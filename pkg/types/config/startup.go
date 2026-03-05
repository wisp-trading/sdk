package config

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// StartupConfig contains everything needed to start a strategy
type StartupConfig struct {
	// Strategy is the loaded strategy configuration
	Strategy *Strategy

	// ConnectorConfigs are the initialized connector configurations
	ConnectorConfigs map[connector.ExchangeName]connector.Config

	// Assets maps each exchange to the pairs declared in config.
	// The runtime routes these to the correct domain watchlist after
	// connector types are known.
	Assets map[connector.ExchangeName][]portfolio.Pair

	// PluginPath is the path to the .so file (for plugin mode)
	PluginPath string

	// StrategyDir is the directory containing the strategy
	StrategyDir string
}

// StartupConfigLoader loads all configuration needed to run a strategy
type StartupConfigLoader interface {
	// LoadForStrategy loads strategy config, connector configs, and asset configs
	// strategyDir: path to the strategy directory containing config.yml
	LoadForStrategy(strategyDir string, wispPath string) (*StartupConfig, error)
}
