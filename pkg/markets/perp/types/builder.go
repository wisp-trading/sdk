package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PerpSignalBuilder provides a fluent API for constructing perpetual futures signals.
// Obtain one via wisp.Perp().Signal(strategyName).
type PerpSignalBuilder interface {
	Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) PerpSignalBuilder
	Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) PerpSignalBuilder
	Build() (PerpSignal, error)
}
