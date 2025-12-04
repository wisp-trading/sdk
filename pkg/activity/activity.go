package activity

import (
	storeActivity "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
)

// activity provides read-only access to positions, trades, and PNL
type activity struct {
	positions kronosActivity.Positions
	trades    kronosActivity.Trades
	pnl       kronosActivity.PNL
}

// NewActivity creates a new Activity instance from the underlying stores
func NewActivity(positionStore storeActivity.Positions, tradeStore storeActivity.Trades) kronosActivity.Activity {
	pos := NewPositions(positionStore)
	trd := NewTrades(tradeStore)
	pnl := NewPNL(pos, trd)

	return &activity{
		positions: pos,
		trades:    trd,
		pnl:       pnl,
	}
}

// Positions returns read-only access to position data
func (a *activity) Positions() kronosActivity.Positions {
	return a.positions
}

// Trades returns read-only access to trade data
func (a *activity) Trades() kronosActivity.Trades {
	return a.trades
}

// PNL returns access to PNL calculations
func (a *activity) PNL() kronosActivity.PNL {
	return a.pnl
}

var _ kronosActivity.Activity = (*activity)(nil)
