package types

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SpotAction represents a single order action for the spot market.
type SpotAction struct {
	strategy.BaseAction
	Pair     portfolio.Pair    `json:"pair"`
	Quantity numerical.Decimal `json:"quantity"`
	Price    numerical.Decimal `json:"price"`
}

// GetMarketType returns the spot market type.
func (a *SpotAction) GetMarketType() connector.MarketType {
	return connector.MarketTypeSpot
}

// Validate checks that the action is well-formed.
// A zero price is treated as a market order and is permitted.
func (a *SpotAction) Validate() error {
	if err := a.ValidateBase(); err != nil {
		return err
	}
	if !a.Pair.Base().IsValid() || !a.Pair.Quote().IsValid() {
		return fmt.Errorf("pair must have valid base and quote")
	}
	if a.Quantity.IsZero() || a.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}
	if a.Price.IsNegative() {
		return fmt.Errorf("price must not be negative")
	}
	return nil
}
