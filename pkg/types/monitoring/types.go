package monitoring

import (
	"time"

	"github.com/wisp-trading/wisp/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

// PnLView represents the PnL snapshot for a strategy
type PnLView struct {
	StrategyName  string            `json:"strategy_name"`
	RealizedPnL   numerical.Decimal `json:"realized_pnl"`
	UnrealizedPnL numerical.Decimal `json:"unrealized_pnl"`
	TotalPnL      numerical.Decimal `json:"total_pnl"`
	TotalFees     numerical.Decimal `json:"total_fees"`
}

// StrategyMetrics represents runtime metrics for a strategy
type StrategyMetrics struct {
	StrategyName     string        `json:"strategy_name"`
	Status           string        `json:"status"`
	LastSignalTime   time.Time     `json:"last_signal_time"`
	SignalsGenerated int           `json:"signals_generated"`
	SignalsExecuted  int           `json:"signals_executed"`
	SignalsFailed    int           `json:"signals_failed"`
	AverageLatency   time.Duration `json:"average_latency"`
	ActivePositions  int           `json:"active_positions"`
	DailyPnL         float64       `json:"daily_pnl"`
	WeeklyPnL        float64       `json:"weekly_pnl"`
	MonthlyPnL       float64       `json:"monthly_pnl"`
}

// Type aliases for SDK profiling types
// These allow us to use SDK profiling types directly in the CLI
type (
	ProfilingStats           = profiling.StrategyStats
	ProfilingMetrics         = profiling.StrategyMetrics
	ProfilingIndicatorTiming = profiling.IndicatorTiming
)
