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

func (b *perpBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

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

func (b *perpBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *perpBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.PerpSignalBuilder {
	b.actions = append(b.actions, &strategy.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

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

func (b *perpBuilder) Build() strategy.PerpSignal {
	return strategy.NewPerpSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions)
}
