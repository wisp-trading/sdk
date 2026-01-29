package signal

import (
	"github.com/google/uuid"
	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
	"github.com/wisp-trading/wisp/pkg/types/strategy"
	"github.com/wisp-trading/wisp/pkg/types/temporal"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

// builder is the concrete implementation of strategy.SignalBuilder.
type builder struct {
	strategyName strategy.StrategyName
	actions      []strategy.TradeAction
	timeProvider temporal.TimeProvider
}

// Buy adds a buy action to the signal.
func (b *builder) Buy(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionBuy,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    numerical.NewFromInt(0),
	})
	return b
}

// BuyLimit adds a limit buy action to the signal.
func (b *builder) BuyLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionBuy,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    price,
	})
	return b
}

// Sell adds a sell action to the signal.
func (b *builder) Sell(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSell,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    numerical.NewFromInt(0),
	})
	return b
}

// SellLimit adds a limit sell action to the signal.
func (b *builder) SellLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSell,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    price,
	})
	return b
}

// SellShort adds a short sell action to the signal.
func (b *builder) SellShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSellShort,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    numerical.NewFromInt(0),
	})
	return b
}

// SellShortLimit adds a limit short sell action to the signal.
func (b *builder) SellShortLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price numerical.Decimal) strategy.SignalBuilder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSellShort,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    price,
	})
	return b
}

// Build constructs the final Signal object with the correct timestamp from the time provider.
func (b *builder) Build() *strategy.Signal {
	return &strategy.Signal{
		ID:        uuid.New(),
		Strategy:  b.strategyName,
		Actions:   b.actions,
		Timestamp: b.timeProvider.Now(),
	}
}
