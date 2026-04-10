package activity

import "github.com/wisp-trading/sdk/pkg/types/execution"

// Activity provides read-only access to positions, trades, PNL, and execution history.
type Activity interface {
	PNL() PNL

	// Executions returns the execution records store, which holds the outcome of
	// every signal dispatched through this session — success, failure, and hook errors.
	Executions() execution.ExecutionRecords
}
