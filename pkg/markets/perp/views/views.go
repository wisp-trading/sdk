package views

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type perpViews struct {
	watchlist         perpTypes.PerpWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewPerpViews(
	watchlist perpTypes.PerpWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) perpTypes.PerpViews {
	return &perpViews{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// GetMarketViews returns all perp markets currently on the watchlist.
// Driven live — always reflects the current state of the perp watchlist.
func (v *perpViews) GetMarketViews() []monitoring.PerpMarketView {
	perpConnectors := v.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().Build())
	result := make([]monitoring.PerpMarketView, 0)

	for _, conn := range perpConnectors {
		info := conn.GetConnectorInfo()
		for _, pair := range v.watchlist.GetRequiredPairs(info.Name) {
			result = append(result, monitoring.PerpMarketView{
				Exchange: string(info.Name),
				Pair:     pair.Symbol(),
			})
		}
	}

	return result
}
