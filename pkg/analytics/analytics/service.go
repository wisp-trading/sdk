package analytics

import (
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// analytics provides user-friendly methods for market analytics.
// It delegates market data access to the Market service which handles spot/perp routing.
type analytics struct {
	market analyticsTypes.Market
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(market analyticsTypes.Market) analyticsTypes.Analytics {
	return &analytics{market: market}
}
