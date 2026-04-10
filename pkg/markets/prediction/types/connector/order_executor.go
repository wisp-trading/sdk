package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// LimitOrder represents a limit order in a blockchain based prediction market.
type LimitOrder struct {
	Outcome     Outcome               `json:"outcome,required"`
	Price       numerical.Decimal
	Amount      numerical.Decimal
	Side        connector.OrderSide   `json:"side,required"`
	Expiration  int64                 `json:"expiration"`
	TimeInForce connector.TimeInForce `json:"time_in_force"` // GTC (default) | FOK | FAK
}

type OrderExecutor interface {
	PlaceLimitOrder(order LimitOrder) (*connector.OrderResponse, error)
	PlaceLimitOrders(orders []LimitOrder) ([]*connector.OrderResponse, error)
}
