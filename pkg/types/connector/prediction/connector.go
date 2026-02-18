package prediction

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// Connector represents a prediction market exchange connection
// Follows the same pattern as spot.Connector and perp.Connector
// Uses standard connector.OrderExecutor and connector.AccountReader for trading
type Connector interface {
	connector.Connector
	OrderExecutor
	AccountReader

	GetMarket(slug string) (Market, error)
	GetRecurringMarket(slug string, interval RecurrenceInterval) (Market, error)
	GetOutcome(marketID, outcomeID string) Outcome
}
