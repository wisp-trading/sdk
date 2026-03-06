package types

import (
	market2 "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
)

// MarketStore handles spot market data storage
// Extends base MarketStore and implements BaseMarketStore for registry compatibility
type MarketStore interface {
	market2.MarketStore
	market2.OrderBookStoreExtension
	market2.KlineStoreExtension
}
