package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/common"
)

// MarketStore handles spot market data storage
type MarketStore interface {
	common.MarketStore
}
