package hooks

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/executor"
	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
)

// LoggingHook logs execution events
type LoggingHook struct {
	logger logging.TradingLogger
}

// NewLoggingHook creates a new logging hook
func NewLoggingHook(logger logging.TradingLogger) *LoggingHook {
	return &LoggingHook{
		logger: logger,
	}
}

// BeforeExecute logs before execution
func (h *LoggingHook) BeforeExecute(ctx *executor.ExecutionContext) error {
	h.logger.Info("🔄 About to execute signal %s with %d actions",
		ctx.Signal.ID, len(ctx.Signal.Actions))
	return nil
}

// AfterExecute logs after successful execution
func (h *LoggingHook) AfterExecute(ctx *executor.ExecutionContext, result *executor.ExecutionResult) error {
	if result.Success {
		h.logger.Info("✅ Signal %s executed successfully. Orders: %v",
			ctx.Signal.ID, result.OrderIDs)
	}
	return nil
}

// OnError logs errors
func (h *LoggingHook) OnError(ctx *executor.ExecutionContext, err error) error {
	h.logger.Failed(
		string(ctx.Signal.Strategy),
		ctx.Signal.Actions[0].Asset.Symbol(),
		"❌ Error executing signal %s: %v",
		ctx.Signal.ID,
		err,
	)
	return err
}

// MetricsHook tracks execution metrics
type MetricsHook struct {
	TotalExecutions int
	TotalSuccess    int
	TotalFailures   int
	TotalOrders     int
}

// NewMetricsHook creates a new metrics hook
func NewMetricsHook() execution.ExecutionHook {
	return &MetricsHook{}
}

// BeforeExecute increments execution counter
func (h *MetricsHook) BeforeExecute(ctx *executor.ExecutionContext) error {
	h.TotalExecutions++
	return nil
}

// AfterExecute tracks success metrics
func (h *MetricsHook) AfterExecute(ctx *executor.ExecutionContext, result *executor.ExecutionResult) error {
	if result.Success {
		h.TotalSuccess++
		h.TotalOrders += len(result.OrderIDs)
	}
	return nil
}

// OnError tracks failure metrics
func (h *MetricsHook) OnError(ctx *executor.ExecutionContext, err error) error {
	h.TotalFailures++
	return err
}

// GetStats returns execution statistics
func (h *MetricsHook) GetStats() string {
	successRate := 0.0
	if h.TotalExecutions > 0 {
		successRate = float64(h.TotalSuccess) / float64(h.TotalExecutions) * 100
	}
	return fmt.Sprintf("Executions: %d | Success: %d (%.1f%%) | Failures: %d | Orders: %d",
		h.TotalExecutions, h.TotalSuccess, successRate, h.TotalFailures, h.TotalOrders)
}
