package profiling

import (
	"container/ring"
	"sort"
	"sync"
	"time"

	profiling2 "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// store implements profiling.ProfilingStore
// Stores metrics in circular buffers per strategy
type store struct {
	// Strategy name -> circular buffer of metrics
	strategyMetrics map[string]*ring.Ring

	mu           sync.RWMutex
	capacity     int // per-strategy buffer size
	timeProvider temporal.TimeProvider

	// Async recording
	recordChan chan profiling2.StrategyMetrics
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewStore creates a new profiling store
func NewStore(capacity int, timeProvider temporal.TimeProvider) profiling2.ProfilingStore {
	store := &store{
		strategyMetrics: make(map[string]*ring.Ring),
		capacity:        capacity,
		timeProvider:    timeProvider,
		recordChan:      make(chan profiling2.StrategyMetrics, 100),
		stopChan:        make(chan struct{}),
	}

	// Start background recorder
	store.wg.Add(1)
	go store.recordLoop()

	return store
}

// NewContext creates a new profiling context for a strategy execution
func (s *store) NewContext(strategyName string) profiling2.Context {
	return newExecutionContext(strategyName, s.timeProvider)
}

// RecordExecution records a completed strategy execution (async)
func (s *store) RecordExecution(metrics profiling2.StrategyMetrics) {
	select {
	case s.recordChan <- metrics:
	case <-s.stopChan:
		// Store is stopping, discard metric
	default:
		// Channel full, discard oldest (non-blocking)
	}
}

// recordLoop processes metrics asynchronously
func (s *store) recordLoop() {
	defer s.wg.Done()

	for {
		select {
		case metrics := <-s.recordChan:
			s.record(metrics)
		case <-s.stopChan:
			// Drain remaining metrics
			for {
				select {
				case metrics := <-s.recordChan:
					s.record(metrics)
				default:
					return
				}
			}
		}
	}
}

// record actually stores the metric
func (s *store) record(metrics profiling2.StrategyMetrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get or create ring buffer for this strategy
	buffer, exists := s.strategyMetrics[metrics.StrategyName]
	if !exists {
		buffer = ring.New(s.capacity)
		s.strategyMetrics[metrics.StrategyName] = buffer
	}

	// Store metric and advance
	buffer.Value = metrics
	s.strategyMetrics[metrics.StrategyName] = buffer.Next()
}

// GetRecentMetrics returns the N most recent metrics for a strategy
func (s *store) GetRecentMetrics(strategyName string, limit int) []profiling2.StrategyMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buffer, exists := s.strategyMetrics[strategyName]
	if !exists {
		return []profiling2.StrategyMetrics{}
	}

	var metrics []profiling2.StrategyMetrics
	count := 0

	// Walk backwards through ring to get most recent
	buffer.Do(func(v interface{}) {
		if v != nil && count < limit {
			if m, ok := v.(profiling2.StrategyMetrics); ok {
				metrics = append(metrics, m)
				count++
			}
		}
	})

	// Reverse to get chronological order
	for i := 0; i < len(metrics)/2; i++ {
		metrics[i], metrics[len(metrics)-1-i] = metrics[len(metrics)-1-i], metrics[i]
	}

	return metrics
}

// GetAverageExecutionTime returns the average execution time for a strategy
func (s *store) GetAverageExecutionTime(strategyName string) time.Duration {
	metrics := s.GetRecentMetrics(strategyName, s.capacity)
	if len(metrics) == 0 {
		return 0
	}

	var total time.Duration
	for _, m := range metrics {
		total += m.ExecutionTime
	}

	return total / time.Duration(len(metrics))
}

// GetPercentile returns the Nth percentile execution time for a strategy
func (s *store) GetPercentile(strategyName string, percentile float64) time.Duration {
	metrics := s.GetRecentMetrics(strategyName, s.capacity)
	if len(metrics) == 0 {
		return 0
	}

	// Extract durations and sort
	durations := make([]time.Duration, len(metrics))
	for i, m := range metrics {
		durations[i] = m.ExecutionTime
	}
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	// Calculate percentile index
	index := int(float64(len(durations)-1) * percentile / 100.0)
	if index < 0 {
		index = 0
	}
	if index >= len(durations) {
		index = len(durations) - 1
	}

	return durations[index]
}

// GetStats returns statistical analysis of strategy performance
func (s *store) GetStats(strategyName string) profiling2.StrategyStats {
	metrics := s.GetRecentMetrics(strategyName, s.capacity)

	if len(metrics) == 0 {
		return profiling2.StrategyStats{
			StrategyName: strategyName,
		}
	}

	var totalDuration time.Duration
	var minDuration time.Duration
	var maxDuration time.Duration
	successCount := 0
	failureCount := 0
	lastExecution := metrics[0].Timestamp

	for i, m := range metrics {
		totalDuration += m.ExecutionTime

		if i == 0 || m.ExecutionTime < minDuration {
			minDuration = m.ExecutionTime
		}
		if m.ExecutionTime > maxDuration {
			maxDuration = m.ExecutionTime
		}

		if m.Success {
			successCount++
		} else {
			failureCount++
		}

		if m.Timestamp.After(lastExecution) {
			lastExecution = m.Timestamp
		}
	}

	avgDuration := totalDuration / time.Duration(len(metrics))
	successRate := float64(successCount) / float64(len(metrics)) * 100.0

	return profiling2.StrategyStats{
		StrategyName:  strategyName,
		TotalRuns:     len(metrics),
		SuccessCount:  successCount,
		FailureCount:  failureCount,
		SuccessRate:   successRate,
		AvgDuration:   avgDuration,
		MinDuration:   minDuration,
		MaxDuration:   maxDuration,
		P50:           s.GetPercentile(strategyName, 50),
		P95:           s.GetPercentile(strategyName, 95),
		P99:           s.GetPercentile(strategyName, 99),
		LastExecution: lastExecution,
	}
}

// Stop gracefully shuts down the profiling store
func (s *store) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}
