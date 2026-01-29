package activity

import (
	wispActivity "github.com/wisp-trading/wisp/pkg/types/wisp/activity"
)

// activity provides read-only access to positions, trades, and PNL
type activity struct {
	positions wispActivity.Positions
	trades    wispActivity.Trades
	pnl       wispActivity.PNL
}

// NewActivity creates a new Activity instance with injected dependencies
func NewActivity(
	positions wispActivity.Positions,
	trades wispActivity.Trades,
	pnl wispActivity.PNL,
) wispActivity.Activity {
	return &activity{
		positions: positions,
		trades:    trades,
		pnl:       pnl,
	}
}

// Positions returns read-only access to position data
func (a *activity) Positions() wispActivity.Positions {
	return a.positions
}

// Trades returns read-only access to trade data
func (a *activity) Trades() wispActivity.Trades {
	return a.trades
}

// PNL returns access to PNL calculations
func (a *activity) PNL() wispActivity.PNL {
	return a.pnl
}

var _ wispActivity.Activity = (*activity)(nil)
