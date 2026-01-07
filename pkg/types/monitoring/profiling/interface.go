package profiling

import (
	"time"
)

// ProfilingStore is the main interface for storing and retrieving performance metrics
// Phase 1: Strategy metrics only
// Future: Add RecordConnectorMetrics(), RecordFlowMetrics() without breaking changes
type ProfilingStore interface {
	// Context creation
	NewContext(strategyName string) Context

	// Strategy metrics (Phase 1)
	RecordExecution(metrics StrategyMetrics)
	GetRecentMetrics(strategyName string, limit int) []StrategyMetrics
	GetAverageExecutionTime(strategyName string) time.Duration
	GetPercentile(strategyName string, percentile float64) time.Duration
	GetStats(strategyName string) StrategyStats

	// Lifecycle
	Stop()
}

// AnomalyDetector detects performance degradation and anomalies
type AnomalyDetector interface {
	// CheckExecution analyzes a single execution duration and returns an alert if anomalous
	CheckExecution(strategyName string, duration time.Duration) Alert

	// Update internal baseline with new measurement
	UpdateBaseline(strategyName string, duration time.Duration)

	// GetBaseline returns the current baseline duration for a strategy
	GetBaseline(strategyName string) time.Duration

	// Reset clears the baseline for a strategy
	Reset(strategyName string)
}

// Context accumulates metrics during a strategy execution
// This is the lightweight object passed through the execution flow
type Context interface {
	// RecordIndicator records the timing of an indicator calculation
	RecordIndicator(name string, duration time.Duration)

	// RecordSignalGeneration records the time spent generating signals
	RecordSignalGeneration(duration time.Duration)

	// Finalize creates the final metrics snapshot for recording
	Finalize(success bool, err error) StrategyMetrics

	// GetStrategyName returns the name of the strategy being profiled
	GetStrategyName() string
}
