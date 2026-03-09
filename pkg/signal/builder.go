package signal

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spotBuilder is the concrete implementation of strategy.SpotSignalBuilder.
type spotBuilder struct {
	strategyName strategy.StrategyName
	actions      []*strategy.SpotAction
	timeProvider temporal.TimeProvider
}

// Buy adds a buy action to the signal.
func (b *spotBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// BuyLimit adds a limit buy action to the signal.
func (b *spotBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// Sell adds a sell action to the signal.
func (b *spotBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// SellLimit adds a limit sell action to the signal.
func (b *spotBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// SellShort adds a short sell action to the signal.
func (b *spotBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// SellShortLimit adds a limit short sell action to the signal.
func (b *spotBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SpotSignalBuilder {
	b.actions = append(b.actions, &strategy.SpotAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSellShort, Exchange: exchange},
		Pair:       pair,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// Build validates the accumulated actions and constructs the SpotSignal.
// Returns an error if the strategy name is empty, no actions have been added,
// or any action fails validation.
func (b *spotBuilder) Build() (strategy.SpotSignal, error) {
	if b.strategyName == "" {
		return nil, fmt.Errorf("strategy name is required")
	}
	if len(b.actions) == 0 {
		return nil, fmt.Errorf("signal must contain at least one action")
	}
	for i, action := range b.actions {
		if err := action.Validate(); err != nil {
			return nil, fmt.Errorf("action %d is invalid: %w", i, err)
		}
	}
	return strategy.NewSpotSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions), nil
}
