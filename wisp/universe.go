package wisp

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
)

// UniverseProvider computes the spot trading universe.
type UniverseProvider interface {
	Universe() wispTypes.Universe
}

type universeProvider struct {
	marketWatchlist   types.MarketWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewUniverseProvider(marketWatchlist types.MarketWatchlist, connectorRegistry registry.ConnectorRegistry) UniverseProvider {
	return &universeProvider{
		marketWatchlist:   marketWatchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the live spot trading universe.
// Called on demand — always reflects the current state of the connector registry
// and market watchlist.
// Perp markets are not included here; use PerpUniverseProvider (pkg/markets/perp).
// Prediction markets are not included here; use PredictionUniverseProvider.
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

	return wispTypes.Universe{
		Exchanges: exchanges,
		Assets:    assets,
	}
}
