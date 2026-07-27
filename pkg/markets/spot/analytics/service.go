package analytics

import (
	baseAnalytics "github.com/wisp-trading/sdk/pkg/markets/base/analytics"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
)

// service implements spotTypes.SpotMarket via shared pair analytics.
type service struct {
	*baseAnalytics.Service
}

// New creates a new spot analytics service.
func New(store spotTypes.MarketStore) spotTypes.SpotMarket {
	return &service{Service: baseAnalytics.NewService(store)}
}

var _ spotTypes.SpotMarket = (*service)(nil)
