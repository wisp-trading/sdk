package signal

import (
	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// perpBuilder is the concrete implementation of strategy.PerpSignalBuilder.
type perpBuilder struct {
	strategyName strategy.StrategyName
	actions      []*strategy.PerpAction
	timeProvider temporal.TimeProvider
}

// Buy adds a market buy action for a perpetual futures position.
func (b *perpBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// BuyLimit adds a limit buy action for a perpetual futures position.
func (b *perpBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// BuyLimitWithLeverage adds a limit buy action with specified leverage.
func (b *perpBuilder) BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
		Leverage:   leverage,
	})
	return b
}

// Sell adds a market sell action to close a perpetual futures position.
func (b *perpBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// SellLimit adds a limit sell action to close a perpetual futures position.
func (b *perpBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// SellShort adds a market short sell action for a perpetual futures position.
func (b *perpBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// SellShortLimit adds a limit short sell action for a perpetual futures position.
func (b *perpBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// SellShortLimitWithLeverage adds a limit short sell action with specified leverage.
func (b *perpBuilder) SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
		Leverage:   leverage,
	})
	return b
}

// Build constructs the final Signal object.
func (b *perpBuilder) Build() *strategy.PerpSignal {
	return &strategy.PerpSignal{
		ID:        uuid.New(),
		Strategy:  b.strategyName,
		Actions:   b.actions,
		Timestamp: b.timeProvider.Now(),
	}
}
