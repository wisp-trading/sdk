package spot

import (
	baseUniverse "github.com/wisp-trading/sdk/pkg/markets/base/universe"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type universeProvider struct {
	watchlist         spotTypes.SpotWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewSpotUniverseProvider(
	watchlist spotTypes.SpotWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) spotTypes.SpotUniverseProvider {
	return &universeProvider{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

func (u *universeProvider) Universe() spotTypes.SpotUniverse {
	ready := u.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().Build())
	names := make([]connector.ExchangeName, 0, len(ready))
	for _, conn := range ready {
		names = append(names, conn.GetConnectorInfo().Name)
	}
	uni := baseUniverse.BuildPairUniverse(names, connector.MarketTypeSpot, u.watchlist)
	return spotTypes.SpotUniverse{
		Exchanges: uni.Exchanges,
		Assets:    uni.Assets,
	}
}
