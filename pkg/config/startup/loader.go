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

	if len(stratConfig.Exchanges) == 0 {
		return nil, fmt.Errorf(
			"strategy config (%s): exchanges is empty — list at least one exchange (e.g. hyperliquid)",
			configPath,
		)
	}

	// 3) map + validate keys for strategy exchanges (hard fail on any problem)
	connectorConfigs, err := l.connectorSvc.GetConnectorConfigsForStrategy(stratConfig.Exchanges)
	if err != nil {
		return nil, fmt.Errorf("connector configs for strategy: %w", err)
	}
	if len(connectorConfigs) == 0 {
		return nil, fmt.Errorf(
			"no connector configs resolved for exchanges %v — Settings → add keys (~/.wisp/connectors.yml)",
			stratConfig.Exchanges,
		)
	}
	l.logger.Info("Resolved connector configs", "count", len(connectorConfigs))

	// 4) pairs (domain routing = connector MarketType)
	assetConfigs := l.convertAssets(stratConfig)
	l.logger.Info("Resolved asset pairs", "count", len(assetConfigs))

	// Helpful mismatch warnings (not hard errors — empty assets can be intentional).
	exchangeSet := make(map[string]struct{}, len(stratConfig.Exchanges))
	for _, ex := range stratConfig.Exchanges {
		exchangeSet[ex] = struct{}{}
	}
	for exName := range stratConfig.Assets {
		if _, ok := exchangeSet[exName]; !ok {
			l.logger.Warn(
				"assets listed for exchange not in exchanges[] — pairs will still load if keys exist",
				"exchange", exName,
			)
		}
	}
	for _, ex := range stratConfig.Exchanges {
		if _, ok := stratConfig.Assets[ex]; !ok {
			l.logger.Warn(
				"exchange listed but has no assets[] entries — domain watchlist starts empty",
				"exchange", ex,
			)
		}
	}

	return &config.StartupConfig{
		Strategy:         stratConfig,
		ConnectorConfigs: connectorConfigs,
		Assets:           assetConfigs,
		StrategyDir:      strategyDir,
	}, nil
}

// convertAssets flattens strategy assets to exchange→pairs.
// Domain loaders keep only exchanges whose connector MarketType matches that domain.
// Skips empty base/quote entries (misconfigured rows).
func (l *startupConfigLoader) convertAssets(
	stratConfig *config.Strategy,
) map[connector.ExchangeName][]portfolio.Pair {
	assets := make(map[connector.ExchangeName][]portfolio.Pair)

	for exName, assetList := range stratConfig.Assets {
		exchange := connector.ExchangeName(exName)
		for _, asset := range assetList {
			if asset.Base == "" || asset.Quote == "" {
				l.logger.Warn(
					"skipping asset with empty base/quote",
					"exchange", exName,
					"base", asset.Base,
					"quote", asset.Quote,
				)
				continue
			}
			pair := portfolio.NewPair(
				portfolio.NewAsset(asset.Base),
				portfolio.NewAsset(asset.Quote),
			)
			assets[exchange] = append(assets[exchange], pair)
		}
	}

	return assets
}
