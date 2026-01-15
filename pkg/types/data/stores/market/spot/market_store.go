package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
)

// MarketStore handles spot market data storage
// Extends base MarketStore and implements BaseMarketStore for registry compatibility
type MarketStore interface {
	market.MarketStore
}
