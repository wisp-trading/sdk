package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// SpotViews owns all monitoring view logic for spot markets.
type SpotViews interface {
	GetMarketViews() []monitoring.SpotMarketView
	GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook
	GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline
}
