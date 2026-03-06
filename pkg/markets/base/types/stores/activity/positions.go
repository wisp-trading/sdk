package activity

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type UpdateKey string
type LastUpdatedMap map[UpdateKey]time.Time

// Positions is the activity store for a single strategy instance.
// Owns orders and trades directly — no strategy key, no StrategyExecution.
type Positions interface {
	// Orders
	AddOrder(order connector.Order)
	UpdateOrderStatus(orderID string, status connector.OrderStatus) error
	GetTotalOrderCount() int64

	// Trades
	AddTrade(trade connector.Trade)
	GetTrades() []connector.Trade

	// Metadata
	GetLastUpdated() LastUpdatedMap
	UpdateLastUpdated(key UpdateKey)
}
