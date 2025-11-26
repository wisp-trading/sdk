package lifecycle

import (
	"sync/atomic"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// TickTimer manages data-driven execution timing with fallback
type TickTimer struct {
	dataUpdateCount      atomic.Int32
	dataUpdatesThreshold int32
	fallbackInterval     time.Duration
	minExecutionInterval time.Duration
	timeProvider         temporal.TimeProvider

	lastExecutionTime atomic.Int64 // Unix nano

	tickChan chan struct{}
	stopChan chan struct{}
	stopped  atomic.Bool
}

// NewTickTimer creates a new tick timer
func NewTickTimer(
	dataUpdatesThreshold int,
	fallbackInterval time.Duration,
	minExecutionInterval time.Duration,
	timeProvider temporal.TimeProvider,
) *TickTimer {
	tt := &TickTimer{
		dataUpdatesThreshold: int32(dataUpdatesThreshold),
		fallbackInterval:     fallbackInterval,
		minExecutionInterval: minExecutionInterval,
		timeProvider:         timeProvider,
		tickChan:             make(chan struct{}, 10),
		stopChan:             make(chan struct{}),
	}

	tt.lastExecutionTime.Store(timeProvider.Now().UnixNano())

	go tt.fallbackLoop()
	return tt
}

// NotifyDataUpdate should be called when new market data arrives
func (t *TickTimer) NotifyDataUpdate() {
	if t.stopped.Load() {
		return
	}

	count := t.dataUpdateCount.Add(1)

	// Atomically check if we hit threshold and reset to 0
	if count >= t.dataUpdatesThreshold {
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

	if now.Sub(lastExec) >= t.minExecutionInterval {
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
	ticker := t.timeProvider.NewTicker(t.fallbackInterval)
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
