package views

import (
	baseViews "github.com/wisp-trading/sdk/pkg/markets/base/views"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type perpViews struct {
	watchlist         perpTypes.PerpWatchlist
	connectorRegistry registry.ConnectorRegistry
	store             perpTypes.MarketStore
}

func NewPerpViews(
	watchlist perpTypes.PerpWatchlist,
	connectorRegistry registry.ConnectorRegistry,
	store perpTypes.MarketStore,
) perpTypes.PerpViews {
	return &perpViews{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
		store:             store,
	}
}

// GetMarketViews returns all perp markets currently on the watchlist.
func (v *perpViews) GetMarketViews() []monitoring.PerpMarketView {
	ready := v.connectorRegistry.FilterPerp(registry.NewFilter().ReadyOnly().Build())
	names := make([]connector.ExchangeName, 0, len(ready))
	for _, conn := range ready {
		names = append(names, conn.GetConnectorInfo().Name)
	}
	refs := baseViews.ListWatchedPairs(names, v.watchlist)
	result := make([]monitoring.PerpMarketView, 0, len(refs))
	for _, r := range refs {
		result = append(result, monitoring.PerpMarketView{Exchange: r.Exchange, Pair: r.Pair})
	}
	return result
}

func (v *perpViews) GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook {
	return baseViews.OrderBook(v.store, exchange, pair)
}

func (v *perpViews) GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	return baseViews.Klines(v.store, exchange, pair, interval, limit)
}
