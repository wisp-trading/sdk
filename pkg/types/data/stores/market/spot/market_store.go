package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
)

// MarketStore handles spot market data storage
type MarketStore interface {
	market.MarketStore
}
