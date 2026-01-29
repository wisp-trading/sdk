package wisp

import (
	"sync"

	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
	"github.com/wisp-trading/wisp/pkg/types/registry"
	wispTypes "github.com/wisp-trading/wisp/pkg/types/wisp"
)

// UniverseProvider computes and caches the trading universe
type UniverseProvider interface {
	Universe() wispTypes.Universe
}

type universeProvider struct {
	assetRegistry     registry.AssetRegistry
	connectorRegistry registry.ConnectorRegistry
	cached            *wispTypes.Universe
	mu                sync.Once
}

// NewUniverseProvider creates a new universe provider with caching
func NewUniverseProvider(assetRegistry registry.AssetRegistry, connectorRegistry registry.ConnectorRegistry) UniverseProvider {
	return &universeProvider{
		assetRegistry:     assetRegistry,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the cached trading universe
func (up *universeProvider) Universe() wispTypes.Universe {
	up.mu.Do(func() {
		// Get ready exchanges
		readyConnectors := up.connectorRegistry.GetAllReadyConnectors()
		exchanges := make([]connector.Exchange, 0, len(readyConnectors))
		for _, conn := range readyConnectors {
			info := conn.GetConnectorInfo()
			exchanges = append(exchanges, connector.Exchange{Name: info.Name})
		}

		// Get assets with their instruments
		assets := make(map[portfolio.Asset][]connector.Instrument)
		requirements := up.assetRegistry.GetAssetRequirements()
		for _, req := range requirements {
			assets[req.Asset] = req.Instruments
		}

		up.cached = &wispTypes.Universe{
			Exchanges: exchanges,
			Assets:    assets,
		}
	})

	return *up.cached
}
