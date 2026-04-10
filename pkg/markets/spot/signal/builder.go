package signal

import (
	"fmt"

	"github.com/google/uuid"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spotBuilder is the concrete implementation of spotTypes.SpotSignalBuilder.
type spotBuilder struct {
	strategyName strategy.StrategyName
	actions      []spotTypes.SpotAction
	timeProvider temporal.TimeProvider
}

func (b *spotBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *spotBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *spotBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *spotBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

func (b *spotBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

func (b *spotBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.actions = append(b.actions, spotTypes.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// Build validates the accumulated actions and constructs the SpotSignal.
func (b *spotBuilder) Build() (spotTypes.SpotSignal, error) {
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
	return spotTypes.NewSpotSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions), nil
}
