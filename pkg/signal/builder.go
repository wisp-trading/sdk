package signal

import (
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

// Build constructs and returns the SpotSignal.
func (b *spotBuilder) Build() strategy.SpotSignal {
	return strategy.NewSpotSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions)
}
