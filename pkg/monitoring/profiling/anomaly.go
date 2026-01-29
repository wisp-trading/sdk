package profiling

import (
	"fmt"
	"sync"
	"time"

	profiling2 "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
)

// anomalyDetector implements profiling.AnomalyDetector
// Uses a moving average baseline with configurable thresholds
type anomalyDetector struct {
	// Strategy name -> baseline duration
	baselines map[string]*baseline

	// Thresholds for alerting
	warningThreshold  float64 // e.g., 1.5x baseline
	criticalThreshold float64 // e.g., 2.0x baseline

	mu sync.RWMutex
}

// baseline tracks the moving average for a strategy
type baseline struct {
	avgDuration time.Duration
	sampleCount int
	maxSamples  int // Window size for moving average
}

// NewAnomalyDetector creates a new anomaly detector
// warningThreshold: multiplier for warning (e.g., 1.5 = 150% of baseline)
// criticalThreshold: multiplier for critical (e.g., 2.0 = 200% of baseline)
// windowSize: number of samples to use for moving average baseline
func NewAnomalyDetector(warningThreshold, criticalThreshold float64, windowSize int) profiling2.AnomalyDetector {
	return &anomalyDetector{
		baselines:         make(map[string]*baseline),
		warningThreshold:  warningThreshold,
		criticalThreshold: criticalThreshold,
	}
}

// CheckExecution analyzes a single execution duration and returns an alert if anomalous
func (d *anomalyDetector) CheckExecution(strategyName string, duration time.Duration) profiling2.Alert {
	d.mu.RLock()
	defer d.mu.RUnlock()

	bl, exists := d.baselines[strategyName]
	if !exists || bl.sampleCount == 0 {
		// No baseline yet
		return profiling2.Alert{
			Severity: profiling2.OK,
			Message:  "No baseline established yet",
		}
	}

	baseline := bl.avgDuration
	ratio := float64(duration) / float64(baseline)

	if ratio >= d.criticalThreshold {
		return profiling2.Alert{
			Severity: profiling2.Critical,
			Message: fmt.Sprintf("Execution time %.2fms is %.1fx baseline (%.2fms) - CRITICAL slowdown",
				float64(duration.Microseconds())/1000.0,
				ratio,
				float64(baseline.Microseconds())/1000.0),
		}
	}

	if ratio >= d.warningThreshold {
		return profiling2.Alert{
			Severity: profiling2.Warning,
			Message: fmt.Sprintf("Execution time %.2fms is %.1fx baseline (%.2fms) - WARNING",
				float64(duration.Microseconds())/1000.0,
				ratio,
				float64(baseline.Microseconds())/1000.0),
		}
	}

	return profiling2.Alert{
		Severity: profiling2.OK,
		Message:  "Execution time within normal range",
	}
}

// UpdateBaseline updates the moving average baseline with a new measurement
func (d *anomalyDetector) UpdateBaseline(strategyName string, duration time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	bl, exists := d.baselines[strategyName]
	if !exists {
		// Initialize new baseline
		d.baselines[strategyName] = &baseline{
			avgDuration: duration,
			sampleCount: 1,
			maxSamples:  100, // Default window size
		}
		return
	}

	// Update moving average
	if bl.sampleCount < bl.maxSamples {
		// Still building up to window size - simple average
		totalDuration := time.Duration(bl.sampleCount) * bl.avgDuration
		totalDuration += duration
		bl.sampleCount++
		bl.avgDuration = totalDuration / time.Duration(bl.sampleCount)
	} else {
		// Window is full - exponential moving average
		// EMA formula: new_avg = alpha * new_value + (1 - alpha) * old_avg
		// alpha = 2 / (window_size + 1)
		alpha := 2.0 / float64(bl.maxSamples+1)
		oldAvg := float64(bl.avgDuration)
		newVal := float64(duration)
		newAvg := alpha*newVal + (1-alpha)*oldAvg
		bl.avgDuration = time.Duration(newAvg)
	}
}

// GetBaseline returns the current baseline duration for a strategy
func (d *anomalyDetector) GetBaseline(strategyName string) time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if bl, exists := d.baselines[strategyName]; exists {
		return bl.avgDuration
	}
	return 0
}

// Reset clears the baseline for a strategy
func (d *anomalyDetector) Reset(strategyName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.baselines, strategyName)
}
