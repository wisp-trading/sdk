package lifecycle

import (
	"context"
	"fmt"

	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type orchestrator struct {
	strategyRegistry registry.StrategyRegistry
	logger           logging.ApplicationLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewOrchestrator creates a new strategy orchestrator.
// The orchestrator is a pure lifecycle manager: it starts and stops strategies.
// Signal routing and execution timing are owned by the strategies themselves.
func NewOrchestrator(
	strategyRegistry registry.StrategyRegistry,
	logger logging.ApplicationLogger,
) lifecycleTypes.Orchestrator {
	return &orchestrator{
		strategyRegistry: strategyRegistry,
		logger:           logger,
	}
}

// Start calls Start on every registered strategy. Each strategy launches its own
// goroutine and manages its own internal clock. Start is non-blocking.
func (o *orchestrator) Start(ctx context.Context) error {
	if o.cancel != nil {
		o.logger.Warn("Orchestrator already started")
		return nil
	}

	o.ctx, o.cancel = context.WithCancel(ctx)

	strategies := o.strategyRegistry.GetAllStrategies()
	o.logger.Info("🎯 Starting strategies", "count", len(strategies))

	for _, strat := range strategies {
		if err := strat.Start(o.ctx); err != nil {
			// Stop any already-started strategies before returning the error.
			o.logger.Error("Failed to start strategy %s: %v — rolling back", strat.GetName(), err)
			_ = o.Stop(ctx)
			return fmt.Errorf("failed to start strategy %s: %w", strat.GetName(), err)
		}
		o.logger.Info("  ✓ Strategy started", "name", strat.GetName())
	}

	o.logger.Info("✅ All strategies started")
	return nil
}

// Stop calls Stop on every registered strategy and waits for them to exit cleanly.
func (o *orchestrator) Stop(ctx context.Context) error {
	if o.cancel == nil {
		return nil
	}

	o.logger.Info("🛑 Stopping strategies")

	// Cancel the orchestrator context — propagated to all strategy contexts.
	o.cancel()
	o.cancel = nil

	for _, strat := range o.strategyRegistry.GetAllStrategies() {
		if err := strat.Stop(ctx); err != nil {
			o.logger.Error("Error stopping strategy %s: %v", strat.GetName(), err)
		}
	}

	o.logger.Info("✅ All strategies stopped")
	return nil
}
