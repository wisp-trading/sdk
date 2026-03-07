package types

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PredictionOrder is the domain-native order record for prediction markets.
// No Pair field — prediction orders are identified by market slug + outcome ID.
type PredictionOrder struct {
	ID          string                        `json:"id"`
	Exchange    connector.ExchangeName        `json:"exchange"`
	MarketID    predictionconnector.MarketID  `json:"market_id"`
	OutcomeID   predictionconnector.OutcomeID `json:"outcome_id"`
	Side        connector.OrderSide           `json:"side"`
	Shares      numerical.Decimal             `json:"shares"`
	Price       numerical.Decimal             `json:"price"`        // probability at placement (0.0–1.0)
	Fee         numerical.Decimal             `json:"fee"`          // fees paid on this order
	RealizedPnL numerical.Decimal             `json:"realized_pnl"` // non-zero once market resolved and redeemed
	Status      connector.OrderStatus         `json:"status"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
}

// PredictionActivityQuery filters prediction orders by exchange and/or market slug.
type PredictionActivityQuery struct {
	Exchange *connector.ExchangeName
	MarketID *predictionconnector.MarketID
}

// PositionsStoreExtension tracks prediction market orders for this instance.
type PositionsStoreExtension interface {
	market.StoreExtension

	// AddOrder records a newly placed prediction order.
	AddOrder(order PredictionOrder)

	// GetOrders returns all orders recorded for this instance.
	GetOrders() []PredictionOrder

	// GetOrdersByExchange returns all orders placed on a specific exchange.
	GetOrdersByExchange(exchange connector.ExchangeName) []PredictionOrder

	// QueryOrders returns orders matching the given query (exchange and/or market slug filter).
	QueryOrders(q PredictionActivityQuery) []PredictionOrder

	// UpdateOrderStatus updates the status of an existing order.
	UpdateOrderStatus(orderID string, status connector.OrderStatus) error

	// UpdateRealizedPnL sets the realized PnL on an order once a market resolves and is redeemed.
	UpdateRealizedPnL(orderID string, realized numerical.Decimal) error
}
