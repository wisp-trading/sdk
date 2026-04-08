package options

import (
	"github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/config"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
)

type assetLoader struct {
	connectorRegistry registryTypes.ConnectorRegistry
	watchlist         types.OptionsWatchlist
	logger            logging.ApplicationLogger
}

// NewAssetLoader creates a new asset loader for options
func NewAssetLoader(
	connectorRegistry registryTypes.ConnectorRegistry,
	watchlist types.OptionsWatchlist,
	logger logging.ApplicationLogger,
) types.OptionsAssetLoader {
	return &assetLoader{
		connectorRegistry: connectorRegistry,
		watchlist:         watchlist,
		logger:            logger,
	}
}

// Load discovers available options expirations across all ready connectors
func (a *assetLoader) Load(cfg *config.StartupConfig) error {
	a.logger.Info("Starting options asset loading")

	// Get all ready options connectors
	optionsConnectors := a.connectorRegistry.FilterOptions(registry.NewFilter().ReadyOnly().Build())
	if len(optionsConnectors) == 0 {
		a.logger.Warn("No ready options connectors found for asset loading")
		return nil
	}

	// For each connector, discover available pairs and their expirations
	for _, conn := range optionsConnectors {
		optionsConn, ok := conn.(optionsconnector.Connector)
		if !ok {
			a.logger.Warnf("Connector %s does not support options interface", conn.GetConnectorInfo().Name)
			continue
		}

		exchangeName := conn.GetConnectorInfo().Name
		a.logger.Debugf("Loading options assets from %s", exchangeName)

		// Get the pairs configured for this exchange
		pairs, exists := cfg.Assets[exchangeName]
		if !exists || len(pairs) == 0 {
			a.logger.Debugf("No pairs configured for %s", exchangeName)
			continue
		}

		// Discover expirations for each pair
		for _, pair := range pairs {
			expirations, err := optionsConn.GetExpirations(pair)
			if err != nil {
				a.logger.Warnf("Failed to get expirations for %s from %s: %v", pair.Symbol(), exchangeName, err)
				continue
			}

			// Add each expiration to the watchlist
			for _, expiration := range expirations {
				if err := a.watchlist.RequireExpiration(exchangeName, pair, expiration); err != nil {
					a.logger.Warnf("Failed to add expiration %s for %s on %s: %v", expiration, pair.Symbol(), exchangeName, err)
				}
			}

			a.logger.Debugf("Loaded %d expirations for %s on %s", len(expirations), pair.Symbol(), exchangeName)
		}

		a.logger.Infof("Completed asset loading for %s", exchangeName)
	}

	a.logger.Info("Options asset loading complete")
	return nil
}
