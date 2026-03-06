package spot

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
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
	readyConnectors := u.connectorRegistry.FilterSpot(
		registry.NewFilter().ReadyOnly().Build(),
	)

	exchanges := make([]connector.Exchange, 0, len(readyConnectors))
	assets := make(map[connector.ExchangeName][]portfolio.Pair)

	for _, conn := range readyConnectors {
		info := conn.GetConnectorInfo()
		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypeSpot,
		})
		if pairs := u.watchlist.GetRequiredPairs(info.Name); len(pairs) > 0 {
			assets[info.Name] = pairs
		}
	}

	return spotTypes.SpotUniverse{
		Exchanges: exchanges,
		Assets:    assets,
	}
}
