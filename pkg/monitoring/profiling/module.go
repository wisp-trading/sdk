package profiling

import (
	"context"

	"github.com/wisp-trading/wisp/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/wisp/pkg/types/temporal"
	"go.uber.org/fx"
)

// Module provides profiling functionality via fx DI
var Module = fx.Module("profiling",
	fx.Provide(
		ProvideProfilingStore,
		ProvideAnomalyDetector,
	),
)

// ProfilingConfig configures the profiling system
type ProfilingConfig struct {
	// Enabled determines if profiling is active
	Enabled bool

	// MetricsCapacity is the number of metrics to keep per strategy
	MetricsCapacity int

	// AnomalyWarningThreshold is the multiplier for warning alerts (e.g., 1.5 = 150% of baseline)
	AnomalyWarningThreshold float64

	// AnomалyCriticalThreshold is the multiplier for critical alerts (e.g., 2.0 = 200% of baseline)
	AnomалyCriticalThreshold float64

	// AnomalyWindowSize is the number of samples for baseline calculation
	AnomalyWindowSize int
}

// DefaultProfilingConfig returns sensible defaults
func DefaultProfilingConfig() ProfilingConfig {
	return ProfilingConfig{
		Enabled:                  true,
		MetricsCapacity:          1000,
		AnomalyWarningThreshold:  1.5,
		AnomалyCriticalThreshold: 2.0,
		AnomalyWindowSize:        100,
	}
}

// ProvideProfilingStore creates a profiling store
func ProvideProfilingStore(
	timeProvider temporal.TimeProvider,
	lc fx.Lifecycle,
) profiling.ProfilingStore {
	config := DefaultProfilingConfig()

	if !config.Enabled {
		return nil // Profiling disabled
	}

	store := NewStore(config.MetricsCapacity, timeProvider)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			store.Stop()
			return nil
		},
	})

	return store
}

// ProvideAnomalyDetector creates an anomaly detector
func ProvideAnomalyDetector() profiling.AnomalyDetector {
	config := DefaultProfilingConfig()

	if !config.Enabled {
		return nil // Profiling disabled
	}

	return NewAnomalyDetector(
		config.AnomalyWarningThreshold,
		config.AnomалyCriticalThreshold,
		config.AnomalyWindowSize,
	)
}
