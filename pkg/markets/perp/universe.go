package perp

import (
	baseUniverse "github.com/wisp-trading/sdk/pkg/markets/base/universe"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type universeProvider struct {
	watchlist         perpTypes.PerpWatchlist
	connectorRegistry registry.ConnectorRegistry
}

// NewPerpUniverseProvider creates a universe provider for the perp domain.
func NewPerpUniverseProvider(
	watchlist perpTypes.PerpWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) perpTypes.PerpUniverseProvider {
	return &universeProvider{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the live perp trading universe — always current, never cached.
func (u *universeProvider) Universe() perpTypes.PerpUniverse {
	ready := u.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().Build())
	names := make([]connector.ExchangeName, 0, len(ready))
	for _, conn := range ready {
		names = append(names, conn.GetConnectorInfo().Name)
	}
	uni := baseUniverse.BuildPairUniverse(names, connector.MarketTypePerp, u.watchlist)
	return perpTypes.PerpUniverse{
		Exchanges: uni.Exchanges,
		Assets:    uni.Assets,
	}
}
