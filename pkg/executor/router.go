package executor

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	profileTypes "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// executorRouter is the default SignalRouter. It dispatches each signal in its own
// goroutine so Route is always non-blocking from the strategy's perspective.
type executorRouter struct {
	executor execution.Executor
	logger   logging.ApplicationLogger
}

// NewExecutorRouter creates a non-profiling signal router backed by an Executor.
func NewExecutorRouter(executor execution.Executor, logger logging.ApplicationLogger) execution.SignalRouter {
	return &executorRouter{executor: executor, logger: logger}
}

func (r *executorRouter) Route(signal strategy.Signal) {
	go func() {
		if err := r.executor.ExecuteSignal(signal); err != nil {
			r.logger.Error("Signal execution failed: strategy=%s id=%s err=%v",
				signal.GetStrategy(), signal.GetID(), err)
		}
	}()
}

// profilingRouter wraps an inner SignalRouter and records per-signal metrics.
// It measures signal age (time from signal build to routing) and execution time,
// and feeds the anomaly detector. Invisible to strategy authors.
type profilingRouter struct {
	inner    execution.SignalRouter
	store    profileTypes.ProfilingStore
	detector profileTypes.AnomalyDetector // may be nil
}

// NewProfilingRouter wraps inner with profiling instrumentation.
func NewProfilingRouter(
	inner execution.SignalRouter,
	store profileTypes.ProfilingStore,
	detector profileTypes.AnomalyDetector, // optional — may be nil
) execution.SignalRouter {
	return &profilingRouter{inner: inner, store: store, detector: detector}
}

func (r *profilingRouter) Route(signal strategy.Signal) {
	signalAge := time.Since(signal.GetTimestamp())

	start := time.Now()
	r.inner.Route(signal)
	executionTime := time.Since(start)

	r.store.RecordExecution(profileTypes.StrategyMetrics{
		StrategyName:  string(signal.GetStrategy()),
		ExecutionTime: executionTime,
		SignalGenTime: signalAge,
		Timestamp:     time.Now(),
		Success:       true,
	})

	if r.detector != nil {
		r.detector.CheckExecution(string(signal.GetStrategy()), executionTime)
		r.detector.UpdateBaseline(string(signal.GetStrategy()), executionTime)
	}
}
