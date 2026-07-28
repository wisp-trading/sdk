package types

import (
	market "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
)

// MarketStore holds onchain trades/orders and optional price snapshots.
type MarketStore interface {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	market.PositionsStoreExtension
}
