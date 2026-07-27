package config

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// StartupConfig is the fully resolved boot blob for StartStandalone.
// Built by StartupConfigLoader.LoadForStrategy:
//
//  1. settingsPath → ResolveSettingsPath → Configuration.LoadSettings (keys)
//  2. strategyDir/config.yml → StrategyConfig.Load (exchanges/assets, no secrets)
//  3. ConnectorService.GetConnectorConfigsForStrategy → validated connector.Config map
//  4. assets flattened for domain watchlists
type StartupConfig struct {
	Strategy *Strategy

	// ConnectorConfigs: SDK-typed configs ready for connector.Initialize.
	ConnectorConfigs map[connector.ExchangeName]connector.Config

	// Assets: exchange → pairs from strategy config.yml.
	// Domain asset loaders keep only exchanges whose connector MarketType matches.
	Assets map[connector.ExchangeName][]portfolio.Pair

	StrategyDir string
}

// StartupConfigLoader assembles StartupConfig for the runtime.
type StartupConfigLoader interface {
	// LoadForStrategy:
	//   strategyDir  — directory with config.yml
	//   settingsPath — optional; empty → ResolveSettingsPath (~/.wisp/connectors.yml)
	LoadForStrategy(strategyDir string, settingsPath string) (*StartupConfig, error)
}
