// Package signal provides the shared pair-market signal builder core.
// Domain packages embed Core and re-export fluent methods with domain return types;
// perp adds leverage methods via AppendWithLeverage.
package signal

import (
	"fmt"

	"github.com/google/uuid"
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Core accumulates pair actions. Domain Build() freezes into spot/perp signal wrappers.
type Core struct {
	StrategyName strategy.StrategyName
	MarketType   connector.MarketType
	Actions      []baseTypes.PairAction
	TimeProvider temporal.TimeProvider
}

// NewCore creates a builder core for the given market type.
func NewCore(
	strategyName strategy.StrategyName,
	timeProvider temporal.TimeProvider,
	marketType connector.MarketType,
) *Core {
	return &Core{
		StrategyName: strategyName,
		MarketType:   marketType,
		Actions:      make([]baseTypes.PairAction, 0),
		TimeProvider: timeProvider,
	}
}

// Append adds a pair action (leverage zero).
func (c *Core) Append(
	actionType strategy.ActionType,
	pair portfolio.Pair,
	exchange connector.ExchangeName,
	quantity, price numerical.Decimal,
) {
	c.Actions = append(c.Actions, baseTypes.NewPairAction(
		actionType, exchange, pair, quantity, price, c.MarketType,
	))
}

// AppendWithLeverage adds a pair action with leverage (perp extension).
func (c *Core) AppendWithLeverage(
	actionType strategy.ActionType,
	pair portfolio.Pair,
	exchange connector.ExchangeName,
	quantity, price, leverage numerical.Decimal,
) {
	c.Actions = append(c.Actions, baseTypes.NewPairAction(
		actionType, exchange, pair, quantity, price, c.MarketType,
	).WithLeverage(leverage))
}

func (c *Core) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) {
	c.Append(strategy.ActionBuy, pair, exchange, quantity, numerical.NewFromInt(0))
}

func (c *Core) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) {
	c.Append(strategy.ActionBuy, pair, exchange, quantity, price)
}

func (c *Core) BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) {
	c.AppendWithLeverage(strategy.ActionBuy, pair, exchange, quantity, price, leverage)
}

func (c *Core) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) {
	c.Append(strategy.ActionSell, pair, exchange, quantity, numerical.NewFromInt(0))
}

func (c *Core) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) {
	c.Append(strategy.ActionSell, pair, exchange, quantity, price)
}

func (c *Core) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) {
	c.Append(strategy.ActionSellShort, pair, exchange, quantity, numerical.NewFromInt(0))
}

func (c *Core) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) {
	c.Append(strategy.ActionSellShort, pair, exchange, quantity, price)
}

func (c *Core) SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) {
	c.AppendWithLeverage(strategy.ActionSellShort, pair, exchange, quantity, price, leverage)
}

// ValidateReady checks strategy name + actions before domain freezes the signal.
func (c *Core) ValidateReady() error {
	if c.StrategyName == "" {
		return fmt.Errorf("strategy name is required")
	}
	if len(c.Actions) == 0 {
		return fmt.Errorf("signal must contain at least one action")
	}
	for i := range c.Actions {
		if err := c.Actions[i].Validate(); err != nil {
			return fmt.Errorf("action %d is invalid: %w", i, err)
		}
	}
	return nil
}

// NewID returns a fresh signal id (domain Build uses this).
func (c *Core) NewID() uuid.UUID {
	return uuid.New()
}
