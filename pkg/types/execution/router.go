package execution

import "github.com/wisp-trading/sdk/pkg/types/strategy"

// SignalRouter routes signals from strategies to the executor.
// The default implementation dispatches each Route call in its own goroutine
// so the calling strategy is never blocked by executor latency.
//
// Strategies must not call this directly — use domain Emit only:
// wisp.Spot().Emit / Perp().Emit / Predict().Emit / Options().Emit.
type SignalRouter interface {
	// Route dispatches the signal fire-and-forget; errors are logged and discarded.
	Route(signal strategy.Signal)

	// RouteWithResult dispatches the signal and sends the ExecutionResult to ch
	// when execution completes (success or failure). The caller must ensure ch
	// is buffered (capacity ≥ 1) to avoid blocking the executor goroutine.
	RouteWithResult(signal strategy.Signal, ch chan<- ExecutionResult)
}

// Dispatch routes a trading signal asynchronously and returns a callback for the result.
// Used by domain Emit helpers (Spot/Perp/Predict/Options).
func Dispatch(router SignalRouter, signal strategy.Signal) ExecutionCallback {
	ch := make(chan ExecutionResult, 1)
	router.RouteWithResult(signal, ch)
	return NewExecutionCallback(ch)
}
