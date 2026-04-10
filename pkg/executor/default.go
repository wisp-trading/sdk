package executor

import (
	"fmt"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
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
	logger          logging.ApplicationLogger
	timeProvider    temporal.TimeProvider
	hookRegistry    registry.Hooks
	execRecords     execution.ExecutionRecords
	spotExecutor    spotTypes.SignalExecutor
	perpExecutor    perpTypes.SignalExecutor
	predExecutor    predTypes.SignalExecutor
	optionsExecutor optionsTypes.SignalExecutor
}

// NewExecutor creates a new default executor
func NewExecutor(
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	hookRegistry registry.Hooks,
	execRecords execution.ExecutionRecords,
	spotExecutor spotTypes.SignalExecutor,
	perpExecutor perpTypes.SignalExecutor,
	predExecutor predTypes.SignalExecutor,
	optionsExecutor optionsTypes.SignalExecutor,
) execution.Executor {
	logger.Info("Initializing executor")
	return &executor{
		logger:          logger,
		timeProvider:    timeProvider,
		hookRegistry:    hookRegistry,
		execRecords:     execRecords,
		spotExecutor:    spotExecutor,
		perpExecutor:    perpExecutor,
		predExecutor:    predExecutor,
		optionsExecutor: optionsExecutor,
	}
}

// ExecuteSignal processes a signal fire-and-forget style; errors are returned but the
// full ExecutionResult is discarded. Use ExecuteSignalWithResult when you need result detail.
func (e *executor) ExecuteSignal(signal strategy.Signal) error {
	result, err := e.ExecuteSignalWithResult(signal)
	if result.HookError != nil {
		e.logger.Error("AfterExecute hook failed for signal %s: %v", signal.GetID(), result.HookError)
	}
	return err
}

// ExecuteSignalWithResult processes a signal and returns the full ExecutionResult.
func (e *executor) ExecuteSignalWithResult(signal strategy.Signal) (execution.ExecutionResult, error) {
	ctx := &execution.ExecutionContext{
		Signal:    signal,
		Timestamp: e.timeProvider.Now(),
		Metadata:  make(map[string]interface{}),
	}

	e.logger.Info("Executing signal %s (strategy: %s)", signal.GetID(), signal.GetStrategy())

	hooks := e.hookRegistry.GetHooks()

	for _, hook := range hooks {
		if err := hook.BeforeExecute(ctx); err != nil {
			e.logger.Warn("Hook blocked execution: %v", err)
			e.handleError(ctx, err, hooks)
			result := execution.ExecutionResult{Success: false, Error: err}
			e.record(signal, result)
			return result, err
		}
	}

	result := execution.ExecutionResult{
		OrderIDs: make([]string, 0),
		Success:  true,
	}

	var execErr error
	switch s := signal.(type) {
	case spotTypes.SpotSignal:
		execErr = e.spotExecutor.ExecuteSpotSignal(s, ctx, &result)
	case perpTypes.PerpSignal:
		execErr = e.perpExecutor.ExecutePerpSignal(s, ctx, &result)
	case predTypes.PredictionSignal:
		execErr = e.predExecutor.ExecutePredictionSignal(s, ctx, &result)
	case optionsTypes.OptionsSignal:
		execErr = e.optionsExecutor.ExecuteOptionsSignal(s, ctx, &result)
	default:
		execErr = fmt.Errorf("unsupported signal type: %T", signal)
	}

	if execErr != nil {
		result.Error = execErr
		result.Success = false
		e.handleError(ctx, execErr, hooks)
		e.record(signal, result)
		return result, execErr
	}

	// Run AfterExecute hooks — failures are recorded in HookError and mark Success=false.
	for _, hook := range hooks {
		if err := hook.AfterExecute(ctx, &result); err != nil {
			e.logger.Error("AfterExecute hook failed: %v", err)
			result.HookError = err
			result.Success = false
		}
	}

	e.record(signal, result)

	if result.Success {
		e.logger.Info("Signal %s executed successfully", signal.GetID())
	}
	return result, nil
}

// record writes the execution outcome to the records store for bookkeeping.
func (e *executor) record(signal strategy.Signal, result execution.ExecutionResult) {
	e.execRecords.Add(execution.ExecutionRecord{
		SignalID:  signal.GetID(),
		Strategy:  signal.GetStrategy(),
		Timestamp: e.timeProvider.Now(),
		OrderIDs:  result.OrderIDs,
		Success:   result.Success,
		Error:     result.Error,
		HookError: result.HookError,
	})
}

// handleError calls OnError hooks
func (e *executor) handleError(ctx *execution.ExecutionContext, err error, hooks []execution.ExecutionHook) {
	for _, hook := range hooks {
		if hookErr := hook.OnError(ctx, err); hookErr != nil {
			e.logger.Error("Hook OnError failed: %v", hookErr)
		}
	}
}

var _ execution.Executor = (*executor)(nil)
