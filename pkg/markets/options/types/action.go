package types

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// OptionsAction represents a single order action for the options market.
type OptionsAction struct {
	strategy.BaseAction
	Contract OptionContract    `json:"contract"`
	Quantity numerical.Decimal `json:"quantity"`
	Price    numerical.Decimal `json:"price"` // Limit price; zero = market order
}

// GetMarketType returns the options market type.
func (a *OptionsAction) GetMarketType() connector.MarketType {
	return connector.MarketTypeOptions
}

// Validate checks that the action is well-formed.
// A zero price is treated as a market order and is permitted.
func (a *OptionsAction) Validate() error {
	if err := a.ValidateBase(); err != nil {
		return err
	}
	if a.Contract.Strike <= 0 {
		return fmt.Errorf("contract strike must be positive")
	}
	if a.Contract.OptionType != "CALL" && a.Contract.OptionType != "PUT" {
		return fmt.Errorf("contract option type must be CALL or PUT")
	}
	if a.Quantity.IsZero() || a.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}
	if a.Price.IsNegative() {
		return fmt.Errorf("price must not be negative")
	}
	return nil
}
