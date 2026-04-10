package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
)

// SignalExecutor executes spot market signals and owns spot trade/order storage.
type SignalExecutor interface {
	ExecuteSpotSignal(
		signal SpotSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error

	// HandleTrade records an inbound trade fill against a pending spot order.
	HandleTrade(trade connector.Trade) error
}
