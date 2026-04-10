package execution

import "time"

// ExecutionCallback is returned by Emit. Strategies that don't care about the
// outcome can discard it; strategies that do can call Await or AwaitWithTimeout.
type ExecutionCallback struct {
	ch <-chan ExecutionResult
}

// NewExecutionCallback wraps a result channel in an ExecutionCallback.
func NewExecutionCallback(ch <-chan ExecutionResult) ExecutionCallback {
	return ExecutionCallback{ch: ch}
}

// Await blocks until the execution completes and returns the result.
func (c ExecutionCallback) Await() ExecutionResult {
	return <-c.ch
}

// AwaitWithTimeout blocks until execution completes or the timeout elapses.
// Returns the result and true on success, or the zero ExecutionResult and false on timeout.
func (c ExecutionCallback) AwaitWithTimeout(d time.Duration) (ExecutionResult, bool) {
	select {
	case r := <-c.ch:
		return r, true
	case <-time.After(d):
		return ExecutionResult{}, false
	}
}
