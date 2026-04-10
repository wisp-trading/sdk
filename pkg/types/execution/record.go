package execution

import (
	"time"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// ExecutionRecord is written after every signal execution — success or failure.
// It is the durable bookkeeping entry for a signal, regardless of whether the
// strategy awaited the ExecutionCallback.
type ExecutionRecord struct {
	SignalID  uuid.UUID
	Strategy  strategy.StrategyName
	Timestamp time.Time
	OrderIDs  []string
	Success   bool
	Error     error // set when the order itself failed
	HookError error // set when the order was placed but an AfterExecute hook failed
}

// ExecutionRecords is the read/write store for execution records.
// The executor writes to it; strategies read via wisp.Activity().Executions().
type ExecutionRecords interface {
	// Add records the outcome of an execution.
	Add(record ExecutionRecord)

	// GetAll returns all execution records in insertion order.
	GetAll() []ExecutionRecord

	// GetBySignalID returns the record for a specific signal, or nil if not found.
	GetBySignalID(id uuid.UUID) *ExecutionRecord

	// GetByStrategy returns all records for a given strategy name.
	GetByStrategy(name strategy.StrategyName) []ExecutionRecord
}
