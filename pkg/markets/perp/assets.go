package perp

import (
	"github.com/wisp-trading/sdk/pkg/markets/perp/types"
	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type perpAssetLoader struct {
	watchlist         types.PerpWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func newPerpAssetLoader(
	watchlist types.PerpWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) types.AssetLoader {
	return &perpAssetLoader{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Load reads the flat assets map from StartupConfig and registers only the pairs
// that belong to perp exchanges, as determined by the connector registry.
func (l *perpAssetLoader) Load(cfg *configTypes.StartupConfig) error {
	for exchange, pairs := range cfg.Assets {
		mt, ok := l.connectorRegistry.ConnectorType(exchange)
		if !ok || mt != connector.MarketTypePerp {
			continue
		}
		for _, pair := range pairs {
			l.watchlist.RequirePair(exchange, pair)
		}
	}
	return nil
}

var _ types.AssetLoader = (*perpAssetLoader)(nil)
