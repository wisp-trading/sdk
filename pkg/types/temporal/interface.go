package temporal

import (
	"time"
)

// TimeProvider abstracts time operations for dual-mode execution
type TimeProvider interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	NewTimer(d time.Duration) Timer
	Since(t time.Time) time.Duration
	NewTicker(d time.Duration) Ticker
	Sleep(d time.Duration)
}

type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

type Ticker interface {
	C() <-chan time.Time
	Reset(d time.Duration)
	Stop()
}
