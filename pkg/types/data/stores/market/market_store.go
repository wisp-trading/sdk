package market

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// StoreExtension allows market-specific data storage (funding rates, etc.)
type StoreExtension interface {
	GetName() string
}

// MarketStore contains shared market data storage methods
type MarketStore interface {
	MarketType() MarketType

	// Order books
	UpdateOrderBook(pair portfolio.Pair, exchange connector.ExchangeName, orderBook connector.OrderBook)
	GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetOrderBooks(pair portfolio.Pair) OrderBookMap
	GetAllPairsWithOrderBooks() []portfolio.Pair

	// Prices
	UpdatePairPrice(pair portfolio.Pair, exchange connector.ExchangeName, price connector.Price)
	UpdatePairPrices(pair portfolio.Pair, prices PriceMap)
	GetPairPrice(pair portfolio.Pair, exchange connector.ExchangeName) *connector.Price
	GetPairPrices(pair portfolio.Pair) PriceMap

	// Klines
	UpdateKline(pair portfolio.Pair, exchange connector.ExchangeName, kline connector.Kline)
	GetKlines(pair portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
	GetKlinesSince(pair portfolio.Pair, exchange connector.ExchangeName, interval string, since time.Time) []connector.Kline

	// Metadata
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)
}

// DataKey represents types of market data (matches old store naming)
type DataKey string

const (
	DataKeyOrderBooks DataKey = "order_books"
	DataKeyPairPrice  DataKey = "pair_price"
	DataKeyKlines     DataKey = "klines"
)

// UpdateKey identifies a specific data update (matches old store structure)
type UpdateKey struct {
	DataType DataKey
	Pair     portfolio.Pair
	Exchange connector.ExchangeName
}

// Type aliases for cleaner return types
type LastUpdatedMap map[UpdateKey]time.Time

type OrderBookMap map[connector.ExchangeName]*connector.OrderBook

type PriceMap map[connector.ExchangeName]connector.Price
type KlineMap map[connector.ExchangeName]map[string][]connector.Kline
