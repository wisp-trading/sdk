package signal

import (
	"fmt"

	"github.com/google/uuid"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// perpBuilder is the concrete implementation of perpTypes.PerpSignalBuilder.
type perpBuilder struct {
	strategyName strategy.StrategyName
	actions      []perpTypes.PerpAction
	timeProvider temporal.TimeProvider
}

func (b *perpBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *perpBuilder) BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
		Leverage:   leverage,
	})
	return b
}

func (b *perpBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *perpBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *perpBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *perpBuilder) SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.actions = append(b.actions, perpTypes.PerpAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
		Leverage:   leverage,
	})
	return b
}

// Build validates the accumulated actions and constructs the PerpSignal.
func (b *perpBuilder) Build() (perpTypes.PerpSignal, error) {
	if b.strategyName == "" {
		return nil, fmt.Errorf("strategy name is required")
	}
	if len(b.actions) == 0 {
		return nil, fmt.Errorf("signal must contain at least one action")
	}
	for i := range b.actions {
		if err := b.actions[i].Validate(); err != nil {
			return nil, fmt.Errorf("action %d is invalid: %w", i, err)
		}
	}
	return perpTypes.NewPerpSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions), nil
}
