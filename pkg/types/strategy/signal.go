package strategy

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SignalFactory creates signal builders.
type SignalFactory interface {
	New(strategyName StrategyName) SignalBuilder
}

// SignalBuilder provides a fluent API for constructing trading signals.
// Each method returns SignalBuilder to enable method chaining.
type SignalBuilder interface {
	Buy(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) SignalBuilder
	BuyLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) SignalBuilder
	Sell(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) SignalBuilder
	SellLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) SignalBuilder
	SellShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal) SignalBuilder
	SellShortLimit(asset portfolio.Asset, exchange connector.ExchangeName, quantity, price decimal.Decimal) SignalBuilder
	Build() *Signal
}

type Signal struct {
	ID        uuid.UUID     `json:"id"`
	Strategy  StrategyName  `json:"strategy"`
	Actions   []TradeAction `json:"actions"`
	Timestamp time.Time     `json:"timestamp"`
}
