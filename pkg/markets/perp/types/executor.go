package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// SignalExecutor is the domain-specific executor interface for perp signals.
// It owns perp order and trade storage.
type SignalExecutor interface {
	ExecutePerpSignal(
		signal strategy.PerpSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error

	// HandleTrade records an inbound trade fill against a pending perp order.
	HandleTrade(trade connector.Trade) error
}
