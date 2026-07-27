package startup

import (
	"fmt"
	"path/filepath"

	"github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type startupConfigLoader struct {
	strategySvc   config.StrategyConfig
	connectorSvc  config.ConnectorService
	configuration config.Configuration
	logger        logging.ApplicationLogger
}

// NewStartupConfigLoader creates a new startup config loader
func NewStartupConfigLoader(
	strategySvc config.StrategyConfig,
	connectorSvc config.ConnectorService,
	configuration config.Configuration,
	logger logging.ApplicationLogger,
) config.StartupConfigLoader {
	return &startupConfigLoader{
		strategySvc:   strategySvc,
		connectorSvc:  connectorSvc,
		configuration: configuration,
		logger:        logger,
	}
}

// LoadForStrategy merges:
//  1. Global keys: settingsPath → ResolveSettingsPath → ~/.wisp/connectors.yml
//  2. Strategy: strategyDir/config.yml (exchanges/assets only)
//  3. Validated connector.Config map for those exchanges
//  4. Asset pairs for domain watchlists
func (l *startupConfigLoader) LoadForStrategy(
	strategyDir string,
	settingsPath string,
) (*config.StartupConfig, error) {
	// 1) credentials
	resolved := config.ResolveSettingsPath(settingsPath)
	_, err := l.configuration.LoadSettings(resolved)
	if err != nil {
		return nil, fmt.Errorf("connector settings (%s): %w", resolved, err)
	}
	l.logger.Info("Loaded connector settings", "path", resolved)

	// 2) strategy (no secrets)
	configPath := filepath.Join(strategyDir, "config.yml")
	stratConfig, err := l.strategySvc.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("strategy config (%s): %w", configPath, err)
	}
	l.logger.Info("Loaded strategy config", "name", stratConfig.Name, "exchanges", stratConfig.Exchanges)

	// 3) map + validate keys for strategy exchanges
	connectorConfigs, err := l.connectorSvc.GetConnectorConfigsForStrategy(stratConfig.Exchanges)
	if err != nil {
		return nil, fmt.Errorf("connector configs for strategy: %w", err)
	}
	l.logger.Info("Resolved connector configs", "count", len(connectorConfigs))

	// 4) pairs (routing to spot/perp/… is by connector market type, not YAML instruments)
	assetConfigs := l.convertAssets(stratConfig)
	l.logger.Info("Resolved asset pairs", "count", len(assetConfigs))

	return &config.StartupConfig{
		Strategy:         stratConfig,
		ConnectorConfigs: connectorConfigs,
		Assets:           assetConfigs,
		StrategyDir:      strategyDir,
	}, nil
}

// convertAssets flattens strategy assets to exchange→pairs.
// Domain loaders keep only exchanges whose connector MarketType matches that domain.
func (l *startupConfigLoader) convertAssets(
	stratConfig *config.Strategy,
) map[connector.ExchangeName][]portfolio.Pair {
	assets := make(map[connector.ExchangeName][]portfolio.Pair)

	for exName, assetList := range stratConfig.Assets {
		exchange := connector.ExchangeName(exName)
		for _, asset := range assetList {
			pair := portfolio.NewPair(
				portfolio.NewAsset(asset.Base),
				portfolio.NewAsset(asset.Quote),
			)
			assets[exchange] = append(assets[exchange], pair)
		}
	}

	return assets
}
