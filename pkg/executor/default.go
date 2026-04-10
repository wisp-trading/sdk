package executor

import (
	"fmt"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// executor routes signals to domain-specific executors.
type executor struct {
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider
	hookRegistry registry.Hooks

	spotExecutor       spotTypes.SignalExecutor
	perpExecutor       perpTypes.SignalExecutor
	predictionExecutor predTypes.SignalExecutor
}

// NewExecutor creates a new default executor
func NewExecutor(
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	hookRegistry registry.Hooks,
	spotExecutor spotTypes.SignalExecutor,
	perpExecutor perpTypes.SignalExecutor,
	predictionExecutor predTypes.SignalExecutor,
) execution.Executor {
	logger.Info("Initializing executor")
	return &executor{
		logger:             logger,
		timeProvider:       timeProvider,
		hookRegistry:       hookRegistry,
		spotExecutor:       spotExecutor,
		perpExecutor:       perpExecutor,
		predictionExecutor: predictionExecutor,
	}
}

// ExecuteSignal processes a signal and executes the associated actions
func (e *executor) ExecuteSignal(signal strategy.Signal) error {
	ctx := &execution.ExecutionContext{
		Signal:    signal,
		Timestamp: e.timeProvider.Now(),
		Metadata:  make(map[string]interface{}),
	}

	e.logger.Info("Executing signal %s (strategy: %s)", signal.GetID(), signal.GetStrategy())

	// Get hooks from registry at execution time
	hooks := e.hookRegistry.GetHooks()

	// Run BeforeExecute hooks
	for _, hook := range hooks {
		if err := hook.BeforeExecute(ctx); err != nil {
			e.logger.Warn("Hook blocked execution: %v", err)
			e.handleError(ctx, err, hooks)
			return err
		}
	}

	// Execute core logic, dispatching on concrete signal type
	result := &execution.ExecutionResult{
		OrderIDs: make([]string, 0),
		Success:  true,
	}

	var execErr error
	switch s := signal.(type) {
	case spotTypes.SpotSignal:
		execErr = e.spotExecutor.ExecuteSpotSignal(s, ctx, result)
	case perpTypes.PerpSignal:
		execErr = e.perpExecutor.ExecutePerpSignal(s, ctx, result)
	case predTypes.PredictionSignal:
		execErr = e.predictionExecutor.ExecutePredictionSignal(s, ctx, result)
	default:
		execErr = fmt.Errorf("unsupported signal type: %T", signal)
	}

	if execErr != nil {
		result.Error = execErr
		result.Success = false
		e.handleError(ctx, execErr, hooks)
		return execErr
	}

	// Run AfterExecute hooks
	for _, hook := range hooks {
		if err := hook.AfterExecute(ctx, result); err != nil {
			e.logger.Error("Hook AfterExecute failed: %v", err)
		}
	}

	e.logger.Info("Signal %s executed successfully", signal.GetID())
	return nil
}

// handleError calls OnError hooks
func (e *executor) handleError(ctx *execution.ExecutionContext, err error, hooks []execution.ExecutionHook) {
	for _, hook := range hooks {
		if hookErr := hook.OnError(ctx, err); hookErr != nil {
			e.logger.Error("Hook OnError failed: %v", hookErr)
		}
	}
}
