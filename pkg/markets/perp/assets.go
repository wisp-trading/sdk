package perp

import (
	baseWatchlist "github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
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

// Load registers pairs for perp-typed exchanges only.
func (l *perpAssetLoader) Load(cfg *configTypes.StartupConfig) error {
	baseWatchlist.LoadPairsForMarketType(cfg, l.watchlist, l.connectorRegistry, connector.MarketTypePerp)
	return nil
}

var _ types.AssetLoader = (*perpAssetLoader)(nil)
