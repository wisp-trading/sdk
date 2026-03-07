package types

import (
	market "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
)

// MarketStore handles spot market data storage.
type MarketStore interface {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	market.PositionsStoreExtension
}
