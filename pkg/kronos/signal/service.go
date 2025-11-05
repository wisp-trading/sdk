package signal

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Service provides signal building functionality with injected time provider.
// This is a singleton injected via fx DI.
type Service struct {
	timeProvider temporal.TimeProvider
}

// NewService creates a new signal service with the injected time provider.
func NewService(timeProvider temporal.TimeProvider) *Service {
	return &Service{
		timeProvider: timeProvider,
	}
}

// New creates a new signal builder for a strategy.
func (s *Service) New(strategyName strategy.StrategyName) *Builder {
	return &Builder{
		strategyName: strategyName,
		actions:      make([]strategy.TradeAction, 0),
		timeProvider: s.timeProvider,
	}
}

// Builder provides a fluent API for constructing trading signals.
type Builder struct {
	strategyName strategy.StrategyName
	actions      []strategy.TradeAction
	timeProvider temporal.TimeProvider
}

// Buy adds a buy action to the signal.
func (b *Builder) Buy(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) *Builder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionBuy,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    decimal.Zero,
	})
	return b
}

// BuyLimit adds a limit buy action to the signal.
func (b *Builder) BuyLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) *Builder {
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
func (b *Builder) Sell(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) *Builder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSell,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    decimal.Zero,
	})
	return b
}

// SellLimit adds a limit sell action to the signal.
func (b *Builder) SellLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) *Builder {
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
func (b *Builder) SellShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) *Builder {
	b.actions = append(b.actions, strategy.TradeAction{
		Action:   strategy.ActionSellShort,
		Asset:    asset,
		Exchange: exchange,
		Quantity: quantity,
		Price:    decimal.Zero,
	})
	return b
}

// SellShortLimit adds a limit short sell action to the signal.
func (b *Builder) SellShortLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) *Builder {
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
func (b *Builder) Build() *strategy.Signal {
	return &strategy.Signal{
		ID:        uuid.New(),
		Strategy:  b.strategyName,
		Actions:   b.actions,
		Timestamp: b.timeProvider.Now(),
	}
}
