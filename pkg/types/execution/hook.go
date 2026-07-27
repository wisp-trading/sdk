package execution

import (
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// Executor is the core interface for executing trading signals.
type Executor interface {
	// ExecuteSignal processes a trading signal fire-and-forget style; returns an error on failure.
	ExecuteSignal(signal strategy.Signal) error

	// ExecuteSignalWithResult processes a trading signal and returns the full ExecutionResult.
	// Use this when the caller needs to inspect order IDs, hook errors, or success status.
	ExecuteSignalWithResult(signal strategy.Signal) (ExecutionResult, error)
}

// ExecutionHook customizes order placement (in-process RegisterHook — not .so plugins).
type ExecutionHook interface {
	// BeforeExecute is called before an order is placed.
	// Return an error to prevent the execution.
	BeforeExecute(ctx *ExecutionContext) error

	// AfterExecute is called after an order is successfully placed.
	AfterExecute(ctx *ExecutionContext, result *ExecutionResult) error

	// OnError is called when an error occurs during execution.
	OnError(ctx *ExecutionContext, err error) error
}
