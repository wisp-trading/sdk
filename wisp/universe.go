package wisp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
)

// UniverseProvider computes the spot/perp trading universe.
// Prediction markets have their own PredictionUniverse under pkg/markets/prediction.
type UniverseProvider interface {
	Universe() wispTypes.Universe
}

type universeProvider struct {
	marketWatchlist   data.MarketWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewUniverseProvider(marketWatchlist data.MarketWatchlist, connectorRegistry registry.ConnectorRegistry) UniverseProvider {
	return &universeProvider{
		marketWatchlist:   marketWatchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the live spot/perp trading universe.
// Called on demand — always reflects the current state of the connector registry
// and market watchlist. Prediction markets are not included here; use PredictionUniverse.
func (up *universeProvider) Universe() wispTypes.Universe {
	assets := make(map[connector.ExchangeName][]portfolio.Pair)
	var exchanges []connector.Exchange

	for _, conn := range up.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().Build()) {
		info := conn.GetConnectorInfo()
		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypeSpot,
		})
		if pairs := up.marketWatchlist.GetRequiredPairs(info.Name); len(pairs) > 0 {
			assets[info.Name] = pairs
		}
	}

	for _, conn := range up.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().Build()) {
		info := conn.GetConnectorInfo()
		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypePerp,
		})
		if pairs := up.marketWatchlist.GetRequiredPairs(info.Name); len(pairs) > 0 {
			assets[info.Name] = pairs
		}
	}

	return wispTypes.Universe{
		Exchanges: exchanges,
		Assets:    assets,
	}
}
