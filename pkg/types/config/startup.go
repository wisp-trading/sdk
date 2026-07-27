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

	// StrategyDir is the directory containing the strategy
	StrategyDir string
}

// StartupConfigLoader loads all configuration needed to run a strategy
type StartupConfigLoader interface {
	// LoadForStrategy loads strategy config.yml + global connector credentials.
	// strategyDir: directory containing config.yml
	// settingsPath: optional explicit path; empty uses ResolveSettingsPath
	// (~/.wisp/connectors.yml, else project-local migration).
	LoadForStrategy(strategyDir string, settingsPath string) (*StartupConfig, error)
}
