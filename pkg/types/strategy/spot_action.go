package strategy

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SpotAction represents an action for spot markets
type SpotAction struct {
	BaseAction
	Pair     portfolio.Pair    `json:"pair"`
	Quantity numerical.Decimal `json:"quantity"`
	Price    numerical.Decimal `json:"price"`
}

// GetMarketType returns spot
func (a *SpotAction) GetMarketType() connector.MarketType {
	return connector.MarketTypeSpot
}

// Validate checks if the spot action is valid
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
	if a.Price.IsZero() || a.Price.IsNegative() {
		return fmt.Errorf("price must be positive")
	}
	return nil
}
