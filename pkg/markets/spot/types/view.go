package types

import "github.com/wisp-trading/sdk/pkg/types/monitoring"

// SpotViews owns all monitoring view logic for spot markets.
type SpotViews interface {
	GetMarketViews() []monitoring.SpotMarketView
}
