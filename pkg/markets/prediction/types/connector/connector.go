package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// MarketsFilter specifies criteria for filtering markets returned by Markets().
type MarketsFilter struct {
	// MinVolume filters markets by minimum 24h volume (e.g., "1000.00")
	MinVolume string
	// MinLiquidity filters markets by minimum liquidity (e.g., "100.00")
	MinLiquidity string
	// Active filters to only active markets (default: true)
	Active *bool
}

// Connector represents a prediction market exchange connection
// Follows the same pattern as spot.Connector and perp.Connector
// Uses standard connector.OrderExecutor and connector.AccountReader for trading
type Connector interface {
	connector.Connector

	OrderExecutor
	AccountReader

	Redeem(market Market) (string, error)
	GetMarket(slug string) (Market, error)
	GetRecurringMarket(slug string, interval RecurrenceInterval) (Market, error)
	GetOutcome(marketID, outcomeID string) Outcome

	// Markets returns markets available on this exchange, optionally filtered.
	// Pass nil filters to get all active markets.
	Markets(filter *MarketsFilter) ([]Market, error)
}
