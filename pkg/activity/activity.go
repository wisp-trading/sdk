package activity

import (
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
)

// activity provides read-only access to positions, trades, and PNL
type activity struct {
	pnl wispActivity.PNL
}

// NewActivity creates a new Activity instance with injected dependencies
func NewActivity(
	pnl wispActivity.PNL,
) wispActivity.Activity {
	return &activity{
		pnl: pnl,
	}
}

// PNL returns access to PNL calculations
func (a *activity) PNL() wispActivity.PNL {
	return a.pnl
}

var _ wispActivity.Activity = (*activity)(nil)
