package types

import (
	"time"

	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PredictionOrder is the domain-native order record for prediction markets.
// No Pair field — prediction orders are identified by market slug + outcome ID.
type PredictionOrder struct {
	ID           string                        `json:"id"`
	StrategyName strategy.StrategyName         `json:"strategy_name"`
	Exchange     connector.ExchangeName        `json:"exchange"`
	MarketSlug   string                        `json:"market_slug"`
	OutcomeID    predictionconnector.OutcomeID `json:"outcome_id"`
	Side         connector.OrderSide           `json:"side"`
	Shares       numerical.Decimal             `json:"shares"`
	Price        numerical.Decimal             `json:"price"` // probability at placement (0.0–1.0)
	Status       connector.OrderStatus         `json:"status"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// PositionsStoreExtension tracks open prediction market orders per strategy.
// Embedded into the prediction MarketStore alongside OrderBookStoreExtension.
type PositionsStoreExtension interface {
	market.StoreExtension

	// AddOrder records a newly placed prediction order for a strategy.
	AddOrder(strategy strategy.StrategyName, order PredictionOrder)

	// GetOrdersByStrategy returns all orders recorded for the given strategy.
	GetOrdersByStrategy(strategy strategy.StrategyName) []PredictionOrder

	// GetStrategyForOrder resolves which strategy placed the order with the given ID.
	GetStrategyForOrder(orderID string) (strategy.StrategyName, bool)

	// UpdateOrderStatus updates the status of an existing order (e.g. filled, cancelled).
	UpdateOrderStatus(orderID string, status connector.OrderStatus) error
}
