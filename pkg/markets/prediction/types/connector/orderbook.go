package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// OrderBook represents an order book for a specific market outcome
type OrderBook struct {
	connector.OrderBook

	MarketID  MarketID  `json:"market_id"`
	OutcomeID OutcomeID `json:"outcome_id"`
}
