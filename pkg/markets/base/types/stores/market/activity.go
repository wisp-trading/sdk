package market

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// ActivityQuery filters activity data by exchange and/or pair.
// Nil fields mean "no filter" (return all).
type ActivityQuery struct {
	Exchange *connector.ExchangeName
	Pair     *portfolio.Pair
}

// TradesStoreExtension tracks trade execution history for a single domain.
type TradesStoreExtension interface {
	StoreExtension

	AddTrade(trade connector.Trade)
	AddTrades(trades []connector.Trade)

	GetAllTrades() []connector.Trade
	GetTradesByExchange(exchange connector.ExchangeName) []connector.Trade
	GetTradesByPair(pair portfolio.Pair) []connector.Trade
	GetTradesSince(since time.Time) []connector.Trade
	GetTradeByID(tradeID string) *connector.Trade
	TradeExists(tradeID string) bool

	// QueryTrades returns trades matching the given query (exchange and/or pair filter).
	QueryTrades(q ActivityQuery) []connector.Trade

	GetTradeCount() int
	GetTotalVolume(pair portfolio.Pair) numerical.Decimal
}

// PositionsStoreExtension tracks placed orders for a single domain.
type PositionsStoreExtension interface {
	StoreExtension

	AddOrder(order connector.Order)
	UpdateOrderStatus(orderID string, status connector.OrderStatus) error
	GetOrders() []connector.Order
	GetTotalOrderCount() int64

	// QueryOrders returns orders matching the given query (exchange and/or pair filter).
	QueryOrders(q ActivityQuery) []connector.Order
}
