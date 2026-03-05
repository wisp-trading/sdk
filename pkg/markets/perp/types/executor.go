package types

import (
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// SignalExecutor is the domain-specific executor interface for perp signals.
type SignalExecutor interface {
	ExecutePerpSignal(
		signal strategy.PerpSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error
}
