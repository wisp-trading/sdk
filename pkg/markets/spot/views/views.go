package views

import (
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type spotViews struct {
	watchlist         spotTypes.SpotWatchlist
	connectorRegistry registry.ConnectorRegistry
	store             spotTypes.MarketStore
}

func NewSpotViews(
	watchlist spotTypes.SpotWatchlist,
	connectorRegistry registry.ConnectorRegistry,
	store spotTypes.MarketStore,
) spotTypes.SpotViews {
	return &spotViews{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
		store:             store,
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

func (v *spotViews) GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook {
	return v.store.GetOrderBook(pair, exchange)
}

func (v *spotViews) GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return v.store.GetKlines(pair, exchange, interval, limit)
}
