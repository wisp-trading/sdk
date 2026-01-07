package profiling

import (
	"time"
)

// StrategyMetrics represents the performance metrics for a single strategy execution
type StrategyMetrics struct {
	StrategyName     string
	ExecutionTime    time.Duration
	IndicatorMetrics map[string]IndicatorTiming
	SignalGenTime    time.Duration
	Timestamp        time.Time
	Success          bool
	Error            string
}

// IndicatorTiming tracks timing for a specific indicator
type IndicatorTiming struct {
	Name     string
	Duration time.Duration
	Calls    int
}

// StrategyStats provides statistical analysis of strategy performance
type StrategyStats struct {
	StrategyName  string
	TotalRuns     int
	SuccessCount  int
	FailureCount  int
	SuccessRate   float64
	AvgDuration   time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	P50           time.Duration
	P95           time.Duration
	P99           time.Duration
	LastExecution time.Time
}

// AlertSeverity represents the severity level of a performance alert
type AlertSeverity int

const (
	OK AlertSeverity = iota
	Warning
	Critical
)

func (s AlertSeverity) String() string {
	switch s {
	case OK:
		return "OK"
	case Warning:
		return "Warning"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// Alert represents a performance anomaly alert
type Alert struct {
	Severity AlertSeverity
	Message  string
}
