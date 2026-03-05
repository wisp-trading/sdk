package perp

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
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
	readyConnectors := u.connectorRegistry.FilterPerp(
		registry.NewFilter().ReadyOnly().Build(),
	)

	exchanges := make([]connector.Exchange, 0, len(readyConnectors))
	assets := make(map[connector.ExchangeName][]portfolio.Pair)

	for _, conn := range readyConnectors {
		info := conn.GetConnectorInfo()
		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypePerp,
		})
		if pairs := u.watchlist.GetRequiredPairs(info.Name); len(pairs) > 0 {
			assets[info.Name] = pairs
		}
	}

	return perpTypes.PerpUniverse{
		Exchanges: exchanges,
		Assets:    assets,
	}
}
