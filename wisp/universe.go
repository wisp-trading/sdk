package wisp

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	wispTypes "github.com/wisp-trading/sdk/pkg/types/wisp"
)

// UniverseProvider computes and caches the trading universe
type UniverseProvider interface {
	Universe() wispTypes.Universe
}

type universeProvider struct {
	marketWatchlist   data.MarketWatchlist
	connectorRegistry registry.ConnectorRegistry
	cached            *wispTypes.Universe
	mu                sync.Once
}

// NewUniverseProvider creates a new universe provider with caching
func NewUniverseProvider(assetRegistry data.MarketWatchlist, connectorRegistry registry.ConnectorRegistry) UniverseProvider {
	return &universeProvider{
		marketWatchlist:   assetRegistry,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the cached trading universe
func (up *universeProvider) Universe() wispTypes.Universe {
	up.mu.Do(func() {
		// 1) Ready exchanges
		readyConnectors := up.connectorRegistry.Filter(
			registry.NewFilter().ReadyOnly().Build(),
		)

		exchanges := make([]connector.Exchange, 0, len(readyConnectors))
		assets := make(map[connector.ExchangeName][]portfolio.Pair)

		for _, conn := range readyConnectors {
			info := conn.GetConnectorInfo()
			exName := info.Name

			exchanges = append(exchanges, connector.Exchange{Name: exName})

			// 2) Pairs required for this exchange from the watchlist
			pairs := up.marketWatchlist.GetRequiredPairs(exName)
			if len(pairs) > 0 {
				assets[exName] = pairs
			}
		}

		up.cached = &wispTypes.Universe{
			Exchanges: exchanges,
			Assets:    assets,
		}
	})

	return *up.cached
}
