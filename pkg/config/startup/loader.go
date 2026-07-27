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

// LoadForStrategy loads ALL configuration needed to run a strategy.
// settingsPath may be empty — resolved to ~/.wisp/connectors.yml (or migration paths).
func (l *startupConfigLoader) LoadForStrategy(
	strategyDir string,
	settingsPath string,
) (*config.StartupConfig, error) {
	resolved := config.ResolveSettingsPath(settingsPath)
	_, err := l.configuration.LoadSettings(resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to load connector settings: %w", err)
	}

	l.logger.Info("Loaded connector settings", "path", resolved)

	// Load strategy config (per-strategy; no secrets)
	configPath := filepath.Join(strategyDir, "config.yml")
	stratConfig, err := l.strategySvc.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy config: %w", err)
	}

	l.logger.Info("Loaded strategy config", "name", stratConfig.Name, "exchanges", stratConfig.Exchanges)

	connectorConfigs, err := l.connectorSvc.GetConnectorConfigsForStrategy(stratConfig.Exchanges)
	if err != nil {
		return nil, fmt.Errorf("failed to get connector configs: %w", err)
	}

	l.logger.Info("Loaded connector configs", "count", len(connectorConfigs))

	assetConfigs := l.convertAssets(stratConfig)
	l.logger.Info("Loaded asset configs", "count", len(assetConfigs))

	return &config.StartupConfig{
		Strategy:         stratConfig,
		ConnectorConfigs: connectorConfigs,
		Assets:           assetConfigs,
		StrategyDir:      strategyDir,
	}, nil
}

// convertAssets converts strategy config assets to a flat exchange→pairs map.
// Market type routing happens later in the runtime, after connector types are known.
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
