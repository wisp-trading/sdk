package strategy

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// ActionType represents the type of action to take (buy, sell, etc.)
type ActionType string

const (
	ActionBuy       ActionType = "buy"
	ActionSell      ActionType = "sell"
	ActionSellShort ActionType = "sell_short"
	ActionCover     ActionType = "cover"

	ActionHold  ActionType = "hold"
	ActionClose ActionType = "close"
)

// Action is the polymorphic interface for all trading actions
// Each market type implements its own action type
type Action interface {
	// GetType returns the action type (buy, sell, etc.)
	GetType() ActionType

	// GetMarketType returns the market type this action is for
	GetMarketType() market.MarketType

	// GetExchange returns the exchange to execute on
	GetExchange() connector.ExchangeName

	// Validate checks if the action is valid
	Validate() error
}

// BaseAction contains common fields shared across all action types
// Use embedding in concrete action types to inherit this functionality
type BaseAction struct {
	ActionType ActionType             `json:"action"`
	Exchange   connector.ExchangeName `json:"exchange"`
}

// GetType returns the action type
func (a *BaseAction) GetType() ActionType {
	return a.ActionType
}

// GetExchange returns the exchange name
func (a *BaseAction) GetExchange() connector.ExchangeName {
	return a.Exchange
}

// ValidateBase performs common validation for all actions
func (a *BaseAction) ValidateBase() error {
	if a.ActionType == "" {
		return fmt.Errorf("action type is required")
	}
	if a.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}
	return nil
}

// PerpAction represents a perpetual futures market action
type PerpAction struct {
	BaseAction
	Pair     portfolio.Pair    `json:"pair"`
	Quantity numerical.Decimal `json:"quantity"`
	Price    numerical.Decimal `json:"price"`
	Leverage numerical.Decimal `json:"leverage,omitempty"` // Optional leverage (1x if not specified)
}

// GetMarketType returns perpetual
func (a *PerpAction) GetMarketType() market.MarketType {
	return market.MarketTypePerp
}

// Validate checks if the perp action is valid
func (a *PerpAction) Validate() error {
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
	if !a.Leverage.IsZero() && (a.Leverage.LessThan(numerical.NewFromInt(1)) || a.Leverage.GreaterThan(numerical.NewFromInt(125))) {
		return fmt.Errorf("leverage must be between 1x and 125x")
	}
	return nil
}
