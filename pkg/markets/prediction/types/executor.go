package types

import "github.com/wisp-trading/sdk/pkg/types/execution"

// SignalExecutor is the domain-specific executor interface for prediction signals.
// Implemented by pkg/markets/prediction/executor.
type SignalExecutor interface {
	ExecutePredictionSignal(
		signal PredictionSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error
}
