package market

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type OrderBookMap map[connector.ExchangeName]*connector.OrderBook

// OrderBookWriter is a narrower interface for components that only write order books
type OrderBookWriter interface {
	UpdateOrderBook(pair portfolio.Pair, exchange connector.ExchangeName, orderBook connector.OrderBook)
}

// OrderBookStoreExtension provides order book data storage
type OrderBookStoreExtension interface {
	StoreExtension
	UpdateOrderBook(pair portfolio.Pair, exchange connector.ExchangeName, orderBook connector.OrderBook)
	GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetOrderBooks(pair portfolio.Pair) OrderBookMap
	GetAllPairsWithOrderBooks() []portfolio.Pair
}
