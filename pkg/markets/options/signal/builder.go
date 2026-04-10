package signal

import (
	"fmt"

	"github.com/google/uuid"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// optionsBuilder is the concrete implementation of optionsTypes.OptionsSignalBuilder.
type optionsBuilder struct {
	strategyName strategy.StrategyName
	actions      []optionsTypes.OptionsAction
	timeProvider temporal.TimeProvider
}

// Buy adds a market buy action for an options contract.
func (b *optionsBuilder) Buy(contract optionsTypes.OptionContract, exchange connector.ExchangeName, quantity numerical.Decimal) optionsTypes.OptionsSignalBuilder {
	b.actions = append(b.actions, optionsTypes.OptionsAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Contract:   contract,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// BuyLimit adds a limit buy action for an options contract.
func (b *optionsBuilder) BuyLimit(contract optionsTypes.OptionContract, exchange connector.ExchangeName, quantity, price numerical.Decimal) optionsTypes.OptionsSignalBuilder {
	b.actions = append(b.actions, optionsTypes.OptionsAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Contract:   contract,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// Sell adds a market sell action for an options contract.
func (b *optionsBuilder) Sell(contract optionsTypes.OptionContract, exchange connector.ExchangeName, quantity numerical.Decimal) optionsTypes.OptionsSignalBuilder {
	b.actions = append(b.actions, optionsTypes.OptionsAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Contract:   contract,
		Quantity:   quantity,
		Price:      numerical.NewFromInt(0),
	})
	return b
}

// SellLimit adds a limit sell action for an options contract.
func (b *optionsBuilder) SellLimit(contract optionsTypes.OptionContract, exchange connector.ExchangeName, quantity, price numerical.Decimal) optionsTypes.OptionsSignalBuilder {
	b.actions = append(b.actions, optionsTypes.OptionsAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Contract:   contract,
		Quantity:   quantity,
		Price:      price,
	})
	return b
}

// Build validates the accumulated actions and constructs the OptionsSignal.
func (b *optionsBuilder) Build() (optionsTypes.OptionsSignal, error) {
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
	return optionsTypes.NewOptionsSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions), nil
}
