package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market"
	spotTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/spot"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// NewStore creates a new spot market store
func NewStore(timeProvider temporal.TimeProvider) spotTypes.MarketStore {
	return market.NewStore(timeProvider)
}
