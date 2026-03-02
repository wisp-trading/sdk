package config

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// StartupConfig contains everything needed to start a strategy
type StartupConfig struct {
	// Strategy is the loaded strategy configuration
	Strategy *Strategy

	// ConnectorConfigs are the initialized connector configurations
	ConnectorConfigs map[connector.ExchangeName]connector.Config

	// AssetConfigs maps assets to their required instruments
	AssetConfigs map[connector.ExchangeName][]portfolio.Pair

	// PluginPath is the path to the .so file (for plugin mode)
	PluginPath string

	// StrategyDir is the directory containing the strategy
	StrategyDir string

	// ExecutionConfig is the parsed execution timing configuration from config.yml
	// Nil if no execution section is defined (strategy will use defaults)
	ExecutionConfig *strategy.ExecutionConfig
}

// StartupConfigLoader loads all configuration needed to run a strategy
type StartupConfigLoader interface {
	// LoadForStrategy loads strategy config, connector configs, and asset configs
	// strategyDir: path to the strategy directory containing config.yml
	LoadForStrategy(strategyDir string, wispPath string) (*StartupConfig, error)
}
