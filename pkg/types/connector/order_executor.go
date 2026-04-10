package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// OrderExecutor handles order placement and management
type OrderExecutor interface {
	PlaceLimitOrder(pair portfolio.Pair, side OrderSide, quantity, price numerical.Decimal) (*OrderResponse, error)
	PlaceMarketOrder(pair portfolio.Pair, side OrderSide, quantity numerical.Decimal) (*OrderResponse, error)
	CancelOrder(orderID string, pair ...portfolio.Pair) (*CancelResponse, error)
	GetOpenOrders(pair ...portfolio.Pair) ([]Order, error)
	GetOrderStatus(orderID string, pair ...portfolio.Pair) (*Order, error)
}

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

func (s OrderSide) ToString() string {
	return string(s)
}

// TimeInForce controls how long an order remains active and how it fills.
type TimeInForce string

const (
	// TimeInForceGTC — Good Till Cancelled: order rests on the book until filled or manually cancelled.
	TimeInForceGTC TimeInForce = "GTC"
	// TimeInForceFOK — Fill Or Kill: the entire order must fill immediately at the stated price,
	// or it is cancelled in full. No partial fills. Preferred for arb where partial execution
	// would leave a directional position.
	TimeInForceFOK TimeInForce = "FOK"
	// TimeInForceFAK — Fill And Kill (IOC): fills whatever quantity is available immediately,
	// cancels any unfilled remainder.
	TimeInForceFAK TimeInForce = "FAK"
)

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
	OrderID       string            `json:"order_id"`
	ClientOrderID string            `json:"client_order_id,omitempty"`
	Symbol        string            `json:"symbol"`
	Status        OrderStatus       `json:"status"`
	Side          OrderSide         `json:"side"`
	Type          OrderType         `json:"type"`
	Quantity      numerical.Decimal `json:"quantity"`
	Price         numerical.Decimal `json:"price,omitempty"`
	FilledQty     numerical.Decimal `json:"filled_quantity"`
	AvgPrice      numerical.Decimal `json:"average_price,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
}

// Order represents an order on the exchange.
type Order struct {
	ID            string            `json:"id"`
	ClientOrderID string            `json:"client_order_id,omitempty"`
	Pair          portfolio.Pair    `json:"pair"`
	Exchange      ExchangeName      `json:"exchange"`
	Side          OrderSide         `json:"side"`
	Type          OrderType         `json:"type"`
	Status        OrderStatus       `json:"status"`
	Quantity      numerical.Decimal `json:"quantity"`
	Price         numerical.Decimal `json:"price,omitempty"`
	FilledQty     numerical.Decimal `json:"filled_quantity"`
	RemainingQty  numerical.Decimal `json:"remaining_quantity"`
	AvgPrice      numerical.Decimal `json:"average_price,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// CancelResponse represents the response after canceling an order.
type CancelResponse struct {
	OrderID       string         `json:"order_id"`
	ClientOrderID string         `json:"client_order_id,omitempty"`
	Pair          portfolio.Pair `json:"pair"`
	Status        OrderStatus    `json:"status"`
	Timestamp     time.Time      `json:"timestamp"`
}
