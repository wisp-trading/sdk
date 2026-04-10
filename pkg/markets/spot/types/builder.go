package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SpotSignalBuilder provides a fluent API for constructing spot market signals.
// Obtain one via wisp.Spot().Signal(strategyName).
type SpotSignalBuilder interface {
	Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	Build() (SpotSignal, error)
}
