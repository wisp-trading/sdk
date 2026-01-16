package lifecycle

import (
	"sync/atomic"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// TickTimer manages data-driven execution timing with fallback
type TickTimer struct {
	config       TickTimerConfig
	timeProvider temporal.TimeProvider

	dataUpdateCount   atomic.Int32
	lastExecutionTime atomic.Int64 // Unix nano

	tickChan chan struct{}
	stopChan chan struct{}
	stopped  atomic.Bool
}

// NewTickTimer creates a new tick timer with the default configuration
func NewTickTimer(
	dataUpdatesThreshold int,
	fallbackInterval time.Duration,
	minExecutionInterval time.Duration,
	timeProvider temporal.TimeProvider,
) *TickTimer {
	config := TickTimerConfig{
		Mode:                 TickTimerModeDataDriven,
		DataUpdatesThreshold: dataUpdatesThreshold,
		FallbackInterval:     fallbackInterval,
		MinExecutionInterval: minExecutionInterval,
	}
	return NewTickTimerWithConfig(config, timeProvider)
}

// NewTickTimerWithConfig creates a new tick timer with the specified configuration
func NewTickTimerWithConfig(config TickTimerConfig, timeProvider temporal.TimeProvider) *TickTimer {
	tt := &TickTimer{
		config:       config,
		timeProvider: timeProvider,
		tickChan:     make(chan struct{}, 10),
		stopChan:     make(chan struct{}),
	}

	tt.lastExecutionTime.Store(timeProvider.Now().UnixNano())

	// Start the appropriate background loop based on mode
	switch config.Mode {
	case TickTimerModeFixed:
		go tt.fixedIntervalLoop()
	case TickTimerModeDataDriven, TickTimerModeHybrid:
		go tt.fallbackLoop()
	}

	return tt
}

// NotifyDataUpdate should be called when new market data arrives
func (t *TickTimer) NotifyDataUpdate() {
	if t.stopped.Load() {
		return
	}

	// In fixed mode, data updates don't trigger execution
	if t.config.Mode == TickTimerModeFixed {
		return
	}

	count := t.dataUpdateCount.Add(1)

	// Atomically check if we hit threshold and reset to 0
	if count >= int32(t.config.DataUpdatesThreshold) {
		if t.dataUpdateCount.CompareAndSwap(count, 0) {
			// We won the race - trigger execution
			t.tryTriggerExecution()
		}
	}
}

// TickChannel returns the channel that signals when strategies should execute
func (t *TickTimer) TickChannel() <-chan struct{} {
	return t.tickChan
}

// Stop stops the tick timer
func (t *TickTimer) Stop() {
	if t.stopped.CompareAndSwap(false, true) {
		close(t.stopChan)
	}
}

// tryTriggerExecution attempts to trigger execution if minimum interval has passed
func (t *TickTimer) tryTriggerExecution() {
	if t.stopped.Load() {
		return
	}

	now := t.timeProvider.Now()
	lastExec := time.Unix(0, t.lastExecutionTime.Load())

	if now.Sub(lastExec) >= t.config.MinExecutionInterval {
		select {
		case t.tickChan <- struct{}{}:
			t.lastExecutionTime.Store(now.UnixNano())
		default:
			// Channel full, execution already pending
		}
	}
}

// fallbackLoop triggers execution on fallback interval if no data updates
func (t *TickTimer) fallbackLoop() {
	ticker := t.timeProvider.NewTicker(t.config.FallbackInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C():
			t.tryTriggerExecution()
		}
	}
}

// fixedIntervalLoop triggers execution at fixed intervals regardless of data updates
func (t *TickTimer) fixedIntervalLoop() {
	ticker := t.timeProvider.NewTicker(t.config.FixedInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C():
			// In fixed mode, always trigger (respecting min interval)
			t.tryTriggerExecution()
		}
	}
}
