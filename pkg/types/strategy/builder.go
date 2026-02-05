package strategy

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SignalBuilder provides a fluent API for constructing trading signals.
// Each method returns SignalBuilder to enable method chaining.
type SignalBuilder interface {
	Buy(asset portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SignalBuilder
	BuyLimit(asset portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SignalBuilder
	Sell(asset portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SignalBuilder
	SellLimit(asset portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SignalBuilder
	SellShort(asset portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SignalBuilder
	SellShortLimit(asset portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SignalBuilder
	Build() *Signal
}
