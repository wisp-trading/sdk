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
	UpdateOrderBook(asset portfolio.Pair, exchange connector.ExchangeName, orderBook connector.OrderBook)
	GetOrderBook(asset portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetOrderBooks(asset portfolio.Pair) OrderBookMap
	GetAllAssetsWithOrderBooks() []portfolio.Pair

	// Prices
	UpdateAssetPrice(asset portfolio.Pair, exchange connector.ExchangeName, price connector.Price)
	UpdateAssetPrices(asset portfolio.Pair, prices PriceMap)
	GetAssetPrice(asset portfolio.Pair, exchange connector.ExchangeName) *connector.Price
	GetAssetPrices(asset portfolio.Pair) PriceMap

	// Klines
	UpdateKline(asset portfolio.Pair, exchange connector.ExchangeName, kline connector.Kline)
	GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
	GetKlinesSince(asset portfolio.Pair, exchange connector.ExchangeName, interval string, since time.Time) []connector.Kline

	// Metadata
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)
}

// DataKey represents types of market data (matches old store naming)
type DataKey string

const (
	DataKeyOrderBooks DataKey = "order_books"
	DataKeyAssetPrice DataKey = "asset_price"
	DataKeyKlines     DataKey = "klines"
)

// UpdateKey identifies a specific data update (matches old store structure)
type UpdateKey struct {
	DataType DataKey
	Asset    portfolio.Pair
	Exchange connector.ExchangeName
}

// Type aliases for cleaner return types
type LastUpdatedMap map[UpdateKey]time.Time

type OrderBookMap map[connector.ExchangeName]*connector.OrderBook

type PriceMap map[connector.ExchangeName]connector.Price
type KlineMap map[connector.ExchangeName]map[string][]connector.Kline
