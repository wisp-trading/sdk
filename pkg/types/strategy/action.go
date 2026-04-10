package strategy

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// ActionType represents the type of trading action to take.
type ActionType string

const (
	ActionBuy       ActionType = "buy"
	ActionSell      ActionType = "sell"
	ActionSellShort ActionType = "sell_short"
	ActionCover     ActionType = "cover"
	ActionHold      ActionType = "hold"
	ActionClose     ActionType = "close"
)

// Action is the polymorphic interface for all trading actions.
// Each domain implements its own concrete action type.
type Action interface {
	GetType() ActionType
	GetMarketType() connector.MarketType
	GetExchange() connector.ExchangeName
	Validate() error
}

// BaseAction contains common fields shared across all action types.
// Embed this in domain-specific action structs.
type BaseAction struct {
	ActionType ActionType             `json:"action"`
	Exchange   connector.ExchangeName `json:"exchange"`
}

// GetType returns the action type.
func (a *BaseAction) GetType() ActionType { return a.ActionType }

// GetExchange returns the exchange name.
func (a *BaseAction) GetExchange() connector.ExchangeName { return a.Exchange }

// ValidateBase performs validation common to all action types.
func (a *BaseAction) ValidateBase() error {
	if a.ActionType == "" {
		return fmt.Errorf("action type is required")
	}
	if a.Exchange == "" {
		return fmt.Errorf("exchange is required")
	}
	return nil
}
