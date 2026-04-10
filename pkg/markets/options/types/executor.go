package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
)

// SignalExecutor is the domain-specific executor interface for options signals.
// It owns options order and position storage.
type SignalExecutor interface {
	ExecuteOptionsSignal(
		signal OptionsSignal,
		ctx *execution.ExecutionContext,
		result *execution.ExecutionResult,
	) error
}

// OptionsExecutor handles low-level order placement for options contracts.
// Used internally by the SignalExecutor implementation.
type OptionsExecutor interface {
	PlaceOrder(order OptionOrder) (*connector.OrderResponse, error)
	CancelOrder(orderID string, exchange connector.ExchangeName) (*connector.CancelResponse, error)
}

// OptionOrder represents a raw options order for direct placement.
type OptionOrder struct {
	Exchange connector.ExchangeName
	Contract OptionContract
	Side     connector.OrderSide
	Quantity float64
	Price    float64 // Limit price, 0 for market order
	OrderID  string  // Optional: for order tracking
}
