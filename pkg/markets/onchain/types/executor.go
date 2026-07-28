package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
)

// SignalExecutor executes onchain swap signals.
type SignalExecutor interface {
	ExecuteOnchainSignal(
		signal OnchainSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error

	HandleTrade(trade connector.Trade) error
}
