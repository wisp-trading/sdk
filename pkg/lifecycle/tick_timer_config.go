package lifecycle

import "time"

// TickTimerConfig holds configuration for the tick timer
type TickTimerConfig struct {
	// Mode determines how the tick timer operates
	Mode TickTimerMode

	// FixedInterval is the interval for fixed-interval mode
	// Used when Mode == TickTimerModeFixed
	FixedInterval time.Duration

	// DataUpdatesThreshold is the number of data updates before triggering execution
	// Used when Mode == TickTimerModeDataDriven or TickTimerModeHybrid
	// Default: 5
	DataUpdatesThreshold int

	// FallbackInterval is how often to execute if no data updates received
	// Used when Mode == TickTimerModeDataDriven or TickTimerModeHybrid
	// Default: 5s
	FallbackInterval time.Duration

	// MinExecutionInterval is the minimum time between executions
	// Used in all modes to prevent excessive execution
	// Default: 100ms
	MinExecutionInterval time.Duration
}

// TickTimerMode defines how the tick timer triggers execution
type TickTimerMode int

const (
	// TickTimerModeDataDriven triggers execution based on data updates with fallback
	// This is the default mode
	TickTimerModeDataDriven TickTimerMode = iota

	// TickTimerModeFixed triggers execution at fixed intervals
	// Good for longer-term strategies (e.g., every 5 minutes, hourly)
	TickTimerModeFixed

	// TickTimerModeHybrid uses data-driven with a fixed minimum interval
	// Combines responsiveness with predictable minimum frequency
	TickTimerModeHybrid
)

// DefaultTickTimerConfig returns the default tick timer configuration
func DefaultTickTimerConfig() TickTimerConfig {
	return TickTimerConfig{
		Mode:                 TickTimerModeDataDriven,
		DataUpdatesThreshold: 5,
		FallbackInterval:     5 * time.Second,
		MinExecutionInterval: 100 * time.Millisecond,
	}
}

// WithFixedInterval creates a config for fixed-interval execution
func WithFixedInterval(interval time.Duration) TickTimerConfig {
	return TickTimerConfig{
		Mode:                 TickTimerModeFixed,
		FixedInterval:        interval,
		MinExecutionInterval: 100 * time.Millisecond,
	}
}

// WithDataDriven creates a config for data-driven execution
func WithDataDriven(updatesThreshold int, fallbackInterval, minInterval time.Duration) TickTimerConfig {
	return TickTimerConfig{
		Mode:                 TickTimerModeDataDriven,
		DataUpdatesThreshold: updatesThreshold,
		FallbackInterval:     fallbackInterval,
		MinExecutionInterval: minInterval,
	}
}
