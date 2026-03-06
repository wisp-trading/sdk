package market

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// StoreExtension allows market-specific data storage (funding rates, etc.)
type StoreExtension interface {
}

// MarketStore contains minimal core market data storage methods
type MarketStore interface {
	MarketType() connector.MarketType

	// Prices
	UpdatePairPrice(pair portfolio.Pair, exchange connector.ExchangeName, price connector.Price)
	UpdatePairPrices(pair portfolio.Pair, prices PriceMap)
	GetPairPrice(pair portfolio.Pair, exchange connector.ExchangeName) *connector.Price
	GetPairPrices(pair portfolio.Pair) PriceMap

	// Metadata
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)
}

// DataKey represents types of market data
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

type LastUpdatedMap map[UpdateKey]time.Time

type PriceMap map[connector.ExchangeName]connector.Price
