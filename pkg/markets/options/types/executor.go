package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// OptionsExecutor handles order execution for options contracts
type OptionsExecutor interface {
	PlaceOrder(order OptionOrder) (*connector.OrderResponse, error)
	CancelOrder(orderID string) (*connector.CancelResponse, error)
	GetOpenOrders() ([]connector.Order, error)
}

// OptionOrder represents an order for an options contract
type OptionOrder struct {
	Contract OptionContract
	Side     connector.OrderSide
	Quantity float64
	Price    float64 // Limit price, 0 for market order
	OrderID  string  // Optional: for order tracking
}
