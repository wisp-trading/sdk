package connector

import (
	"fmt"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// Market represents a tradeable prediction market
type Market struct {
	// Polymarket: this is the condition id
	MarketID MarketID               `json:"market_id,omitempty"`
	Slug     string                 `json:"slug"`
	Exchange connector.ExchangeName `json:"exchange"`

	// Market structure
	OutcomeType OutcomeType `json:"outcome_type"`
	Outcomes    []Outcome   `json:"outcomes"`

	// Trading status
	Active bool `json:"active"`
	Closed bool `json:"closed"`

	// Timing
	ResolutionTime *time.Time `json:"resolution_time,omitempty"`
	StartTime      *time.Time `json:"start_time,omitempty"`

	RecurringMarket *RecurringMarket `json:"recurring_market,omitempty"`
}

func (m *Market) Validate() error {
	if len(m.Outcomes) == 0 {
		return fmt.Errorf("market must have at least one outcome")
	}

	if m.Slug == "" {
		return fmt.Errorf("market slug must be set")
	}

	for _, outcome := range m.Outcomes {
		if outcome.Pair.Validate() != nil {
			return fmt.Errorf("outcome %s has invalid prediction pair: base and quote must be set", outcome.Pair.Outcome())
		}
	}

	return nil
}

func (m *Market) FindOutcomeById(outcomeId OutcomeID) (*Outcome, error) {
	for _, outcome := range m.Outcomes {
		if outcome.OutcomeID == outcomeId {
			return &outcome, nil
		}
	}

	return nil, fmt.Errorf("no outcome found for outcome id %s", outcomeId)
}

// Outcome represents a tradeable outcome (YES or NO for binary)
type Outcome struct {
	Pair PredictionPair

	// Polymarket: this is the orderbook id
	OutcomeID OutcomeID           `json:"outcome_id,omitempty"`
	Side      connector.OrderSide `json:"side,omitempty"`
}

// OutcomeType represents the market structure
type OutcomeType string

const (
	OutcomeTypeBinary      OutcomeType = "binary"
	OutcomeTypeCategorical OutcomeType = "categorical"
	OutcomeTypeScalar      OutcomeType = "scalar"
)
