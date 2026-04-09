package types

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
)

// MarketStore handles prediction market data storage.
// Embeds base.MarketStore and all prediction-specific extensions.
type MarketStore interface {
	market.MarketStore
	MarketStoreExtension
	OrderBookStoreExtension
	PositionsStoreExtension
	BalanceStoreExtension
}
