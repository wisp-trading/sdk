package connector

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderSide represents the side of an order (buy or sell).
type OrderSide string

const (
	OrderSideBuy     OrderSide = "BUY"
	OrderSideSell    OrderSide = "SELL"
	OrderSideUnknown OrderSide = "UNKNOWN"
)

func FromString(side string) OrderSide {
	switch side {
	case "BUY":
		return OrderSideBuy
	case "SELL":
		return OrderSideSell
	default:
		return OrderSideUnknown
	}
}

func (s OrderSide) IsValid() bool {
	return s == OrderSideBuy || s == OrderSideSell
}

// OrderType represents the type of an order.
type OrderType string

const (
	// Basic order types
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"

	// Stop orders (risk management)
	OrderTypeStopLimit  OrderType = "STOP_LIMIT"
	OrderTypeStopMarket OrderType = "STOP_MARKET"

	// Take profit orders (profit optimization)
	OrderTypeTakeProfitLimit  OrderType = "TAKE_PROFIT_LIMIT"
	OrderTypeTakeProfitMarket OrderType = "TAKE_PROFIT_MARKET"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusOpen            OrderStatus = "OPEN"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderCancellationRequested OrderStatus = "CANCELLATION_REQUESTED"
	OrderStatusPending         OrderStatus = "PENDING"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// OrderResponse represents the response after placing an order.
type OrderResponse struct {
	OrderID       string          `json:"order_id"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
	Symbol        string          `json:"symbol"`
	Status        OrderStatus     `json:"status"`
	Side          OrderSide       `json:"side"`
	Type          OrderType       `json:"type"`
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price,omitempty"`
	FilledQty     decimal.Decimal `json:"filled_quantity"`
	AvgPrice      decimal.Decimal `json:"average_price,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}

// Order represents an order on the exchange.
type Order struct {
	ID            string          `json:"id"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Type          OrderType       `json:"type"`
	Status        OrderStatus     `json:"status"`
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price,omitempty"`
	FilledQty     decimal.Decimal `json:"filled_quantity"`
	RemainingQty  decimal.Decimal `json:"remaining_quantity"`
	AvgPrice      decimal.Decimal `json:"average_price,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// CancelResponse represents the response after canceling an order.
type CancelResponse struct {
	OrderID       string      `json:"order_id"`
	ClientOrderID string      `json:"client_order_id,omitempty"`
	Symbol        string      `json:"symbol"`
	Status        OrderStatus `json:"status"`
	Timestamp     time.Time   `json:"timestamp"`
}
