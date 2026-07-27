package types

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PairAction is the shared order action for pair markets (spot/perp).
// Domains set MarketType at build time; Leverage is the perp extension (zero on spot).
type PairAction struct {
	strategy.BaseAction
	Pair       portfolio.Pair       `json:"pair"`
	Quantity   numerical.Decimal    `json:"quantity"`
	Price      numerical.Decimal    `json:"price"`
	Leverage   numerical.Decimal    `json:"leverage,omitempty"`
	MarketType connector.MarketType `json:"market_type,omitempty"`
}

// GetMarketType returns the domain market type stamped by the builder.
func (a *PairAction) GetMarketType() connector.MarketType {
	if a.MarketType == "" {
		return connector.MarketTypeUnknown
	}
	return a.MarketType
}

// Validate checks base fields + pair/qty/price (shared by spot and perp).
// Zero price is a market order and is permitted. Leverage is not validated here.
func (a *PairAction) Validate() error {
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

// NewPairAction constructs a pair action for the given market type.
func NewPairAction(
	actionType strategy.ActionType,
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	quantity, price numerical.Decimal,
	marketType connector.MarketType,
) PairAction {
	return PairAction{
		BaseAction: strategy.BaseAction{ActionType: actionType, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
		MarketType: marketType,
	}
}

// WithLeverage returns a copy with leverage set (perp extension).
func (a PairAction) WithLeverage(leverage numerical.Decimal) PairAction {
	a.Leverage = leverage
	return a
}
