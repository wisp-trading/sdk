package connector

import (
	"fmt"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
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

	// Category tags for URL construction
	Tags []Tag `json:"tags,omitempty"`

	// EventSlug is the parent event slug used to build the Polymarket URL.
	// A market belongs to an event; the event slug is shorter and is what
	// Polymarket uses in its canonical URLs (e.g. /event/us-x-iran-ceasefire-by).
	EventSlug string `json:"event_slug,omitempty"`
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

// PolymarketURL returns the canonical Polymarket URL for this market.
// It uses the parent event slug (e.g. "us-x-iran-ceasefire-by") which Polymarket
// routes correctly to the event/market page. Falls back to the market slug if no
// event slug is available.
func (m *Market) PolymarketURL() string {
	slug := m.EventSlug
	if slug == "" {
		slug = m.Slug
	}
	return fmt.Sprintf("https://polymarket.com/event/%s", slug)
}

// Tag represents a market category tag
type Tag struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Slug  string `json:"slug"`
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

type Balance struct {
	OutcomeID OutcomeID
	Balance   numerical.Decimal
	Allowance numerical.Decimal
}
