package hooks

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/executor"
	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	"github.com/shopspring/decimal"
)

// RiskHook provides basic risk management
type RiskHook struct {
	MaxPositionSize decimal.Decimal
	MaxDailyTrades  int
	tradesExecuted  int
}

// NewRiskHook creates a new basic risk hook
func NewRiskHook(maxPositionSize decimal.Decimal, maxDailyTrades int) execution.ExecutionHook {
	return &RiskHook{
		MaxPositionSize: maxPositionSize,
		MaxDailyTrades:  maxDailyTrades,
		tradesExecuted:  0,
	}
}

// BeforeExecute checks risk limits before execution
func (h *RiskHook) BeforeExecute(ctx *executor.ExecutionContext) error {
	// Check position size limit
	for _, action := range ctx.Signal.Actions {
		if action.Quantity.GreaterThan(h.MaxPositionSize) {
			return fmt.Errorf("position size %v exceeds max %v for %s",
				action.Quantity, h.MaxPositionSize, action.Asset.Symbol())
		}
	}

	// Check daily trade limit
	if h.MaxDailyTrades > 0 && h.tradesExecuted >= h.MaxDailyTrades {
		return fmt.Errorf("daily trade limit reached (%d/%d)",
			h.tradesExecuted, h.MaxDailyTrades)
	}

	return nil
}

// AfterExecute increments trade counter
func (h *RiskHook) AfterExecute(ctx *executor.ExecutionContext, result *executor.ExecutionResult) error {
	if result.Success {
		h.tradesExecuted++
	}
	return nil
}

// OnError is called when an error occurs
func (h *RiskHook) OnError(ctx *executor.ExecutionContext, err error) error {
	return err
}

// ResetDailyCounter resets the daily trade counter (call this daily)
func (h *RiskHook) ResetDailyCounter() {
	h.tradesExecuted = 0
}
