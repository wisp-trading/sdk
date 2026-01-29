package profiling

import (
	"sync"
	"time"

	profiling2 "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// executionContext implements profiling.Context
// Accumulates metrics during a single strategy execution
type executionContext struct {
	strategyName  string
	startTime     time.Time
	indicators    map[string]*profiling2.IndicatorTiming
	signalGenTime time.Duration
	timeProvider  temporal.TimeProvider
	mu            sync.Mutex
}

// newExecutionContext creates a new profiling context for a strategy execution
func newExecutionContext(strategyName string, timeProvider temporal.TimeProvider) profiling2.Context {
	return &executionContext{
		strategyName: strategyName,
		startTime:    timeProvider.Now(),
		indicators:   make(map[string]*profiling2.IndicatorTiming),
		timeProvider: timeProvider,
	}
}

// RecordIndicator records the timing of an indicator calculation
func (ec *executionContext) RecordIndicator(name string, duration time.Duration) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if timing, exists := ec.indicators[name]; exists {
		// Update existing timing
		timing.Duration += duration
		timing.Calls++
	} else {
		// Create new timing
		ec.indicators[name] = &profiling2.IndicatorTiming{
			Name:     name,
			Duration: duration,
			Calls:    1,
		}
	}
}

// RecordSignalGeneration records the time spent generating signals
func (ec *executionContext) RecordSignalGeneration(duration time.Duration) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.signalGenTime = duration
}

// Finalize creates the final metrics snapshot for recording
func (ec *executionContext) Finalize(success bool, err error) profiling2.StrategyMetrics {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Convert map to map for StrategyMetrics
	indicatorMetrics := make(map[string]profiling2.IndicatorTiming)
	for name, timing := range ec.indicators {
		indicatorMetrics[name] = *timing
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return profiling2.StrategyMetrics{
		StrategyName:     ec.strategyName,
		ExecutionTime:    ec.timeProvider.Since(ec.startTime),
		IndicatorMetrics: indicatorMetrics,
		SignalGenTime:    ec.signalGenTime,
		Timestamp:        ec.startTime,
		Success:          success,
		Error:            errMsg,
	}
}

// GetStrategyName returns the name of the strategy being profiled
func (ec *executionContext) GetStrategyName() string {
	return ec.strategyName
}
