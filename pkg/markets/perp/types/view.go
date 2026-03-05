package types

import "github.com/wisp-trading/sdk/pkg/types/monitoring"

// PerpViews owns all monitoring view logic for perp markets.
// The monitoring ViewRegistry delegates to this interface — it does not implement
// perp-specific logic itself.
type PerpViews interface {
	// GetMarketViews returns all perp markets currently being watched,
	// structured as PerpMarketView entries. Driven live from the perp watchlist.
	GetMarketViews() []monitoring.PerpMarketView
}
