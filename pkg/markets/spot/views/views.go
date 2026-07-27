package views

import (
	baseViews "github.com/wisp-trading/sdk/pkg/markets/base/views"
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
	ready := v.connectorRegistry.FilterSpot(registry.NewFilter().ReadyOnly().Build())
	names := make([]connector.ExchangeName, 0, len(ready))
	for _, conn := range ready {
		names = append(names, conn.GetConnectorInfo().Name)
	}
	refs := baseViews.ListWatchedPairs(names, v.watchlist)
	result := make([]monitoring.SpotMarketView, 0, len(refs))
	for _, r := range refs {
		result = append(result, monitoring.SpotMarketView{Exchange: r.Exchange, Pair: r.Pair})
	}
	return result
}

func (v *spotViews) GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook {
	return baseViews.OrderBook(v.store, exchange, pair)
}

func (v *spotViews) GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return baseViews.Klines(v.store, exchange, pair, interval, limit)
}
