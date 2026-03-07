package spot

import (
	"github.com/wisp-trading/sdk/pkg/markets/spot/types"
	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type spotAssetLoader struct {
	watchlist         types.SpotWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func newSpotAssetLoader(
	watchlist types.SpotWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) types.AssetLoader {
	return &spotAssetLoader{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Load reads the flat assets map from StartupConfig and registers only the pairs
// that belong to spot exchanges, as determined by the connector registry.
func (l *spotAssetLoader) Load(cfg *configTypes.StartupConfig) error {
	for exchange, pairs := range cfg.Assets {
		mt, ok := l.connectorRegistry.ConnectorType(exchange)
		if !ok || mt != connector.MarketTypeSpot {
			continue
		}
		for _, pair := range pairs {
			l.watchlist.RequirePair(exchange, pair)
		}
	}
	return nil
}

var _ types.AssetLoader = (*spotAssetLoader)(nil)
