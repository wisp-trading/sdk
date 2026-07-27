package spot

import (
	baseWatchlist "github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
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

// Load registers pairs for spot-typed exchanges only.
func (l *spotAssetLoader) Load(cfg *configTypes.StartupConfig) error {
	baseWatchlist.LoadPairsForMarketType(cfg, l.watchlist, l.connectorRegistry, connector.MarketTypeSpot)
	return nil
}

var _ types.AssetLoader = (*spotAssetLoader)(nil)
