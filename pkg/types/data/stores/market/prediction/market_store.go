package prediction

import (
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
)

// MarketStore handles prediction market data storage
// Embeds base.MarketStore and adds prediction-specific methods
type MarketStore interface {
	market.MarketStore
}
