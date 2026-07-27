package watchlist

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// LoadPairsForMarketType seeds a pair watchlist from StartupConfig for exchanges
// whose registered MarketType matches want (spot, perp, …).
func LoadPairsForMarketType(
	cfg *configTypes.StartupConfig,
	wl baseTypes.MarketWatchlist,
	connectorRegistry registry.ConnectorRegistry,
	want connector.MarketType,
) {
	if cfg == nil || wl == nil {
		return
	}
	for exchange, pairs := range cfg.Assets {
		mt, ok := connectorRegistry.ConnectorType(exchange)
		if !ok || mt != want {
			continue
		}
		for _, pair := range pairs {
			wl.RequirePair(exchange, pair)
		}
	}
}
