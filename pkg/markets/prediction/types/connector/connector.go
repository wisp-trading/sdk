package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// MarketsFilter specifies criteria for filtering markets returned by Markets().
type MarketsFilter struct {
	// Pagination
	Limit  *int
	Offset *int

	// Volume filters
	MinVolume string
	MaxVolume string

	// Liquidity filters
	MinLiquidity string
	MaxLiquidity string

	// Date range filters (ISO 8601 strings)
	MinEndDate string
	MaxEndDate string

	// Status filter
	Active *bool
	Closed *bool
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

	// FetchOrderBooksForMarket fetches all orderbooks for a market in a single batch request.
	FetchOrderBooksForMarket(market Market) (map[string]*OrderBook, error)
}
