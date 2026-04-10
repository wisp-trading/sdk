package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/execution"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
)

// activity provides read-only access to positions, trades, PNL, and execution history.
type activity struct {
	pnl        wispActivity.PNL
	executions execution.ExecutionRecords
}

// NewActivity creates a new Activity instance with injected dependencies.
func NewActivity(
	pnl wispActivity.PNL,
	executions execution.ExecutionRecords,
) wispActivity.Activity {
	return &activity{
		pnl:        pnl,
		executions: executions,
	}
}

func (a *activity) PNL() wispActivity.PNL {
	return a.pnl
}

func (a *activity) Executions() execution.ExecutionRecords {
	return a.executions
}

var _ wispActivity.Activity = (*activity)(nil)
