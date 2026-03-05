package execution

import "github.com/wisp-trading/sdk/pkg/types/strategy"

// SignalRouter routes signals from strategies to the executor.
// The default implementation dispatches each Route call in its own goroutine
// so the calling strategy is never blocked by executor latency.
type SignalRouter interface {
	Route(signal strategy.Signal)
}
