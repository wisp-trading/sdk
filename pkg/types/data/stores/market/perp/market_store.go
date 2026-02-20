package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
)

// MarketStore handles perpetual market data storage
// Embeds base MarketStore and perp-specific extensions
type MarketStore interface {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	FundingRateStoreExtension
}
