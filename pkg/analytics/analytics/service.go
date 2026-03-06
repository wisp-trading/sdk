package analytics

import (
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
)

// analytics is the concrete implementation of analyticsTypes.Analytics.
// All methods are pure functions — no market or store access occurs here.
type analytics struct{}

// NewAnalyticsService creates a new analytics service.
func NewAnalyticsService() analyticsTypes.Analytics {
	return &analytics{}
}
