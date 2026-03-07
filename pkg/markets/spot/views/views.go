package views

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type spotViews struct {
	watchlist         spotTypes.SpotWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewSpotViews(
	watchlist spotTypes.SpotWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) spotTypes.SpotViews {
	return &spotViews{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

func (v *spotViews) GetMarketViews() []monitoring.SpotMarketView {
	spotConnectors := v.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().Build())
	result := make([]monitoring.SpotMarketView, 0)

	for _, conn := range spotConnectors {
		info := conn.GetConnectorInfo()
		for _, pair := range v.watchlist.GetRequiredPairs(info.Name) {
			result = append(result, monitoring.SpotMarketView{
				Exchange: string(info.Name),
				Pair:     pair.Symbol(),
			})
		}
	}

	return result
}
