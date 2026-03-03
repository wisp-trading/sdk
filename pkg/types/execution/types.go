package execution

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// ExecutionContext contains the context for an execution
type ExecutionContext struct {
	Signal    strategy.Signal
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// ExecutionResult contains the result of an execution
type ExecutionResult struct {
	OrderIDs []string
	Success  bool
	Error    error
}
