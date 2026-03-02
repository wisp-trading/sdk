package strategy

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SpotSignalBuilder provides a fluent API for constructing spot market trading signals.
type SpotSignalBuilder interface {
	Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) SpotSignalBuilder
	SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) SpotSignalBuilder
	Build() *SpotSignal
}

// PerpSignalBuilder provides a fluent API for constructing perpetual futures trading signals.
type PerpSignalBuilder interface {
	Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) PerpSignalBuilder
	Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) PerpSignalBuilder
	SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) PerpSignalBuilder
	SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) PerpSignalBuilder
	Build() *PerpSignal
}

// PredictionSignalBuilder provides a fluent API for constructing prediction market trading signals.
type PredictionSignalBuilder interface {
	Buy(market prediction.Market, outcome prediction.Outcome, exchange connector.ExchangeName, shares, maxPrice numerical.Decimal, expiration int64) PredictionSignalBuilder
	Sell(market prediction.Market, outcome prediction.Outcome, exchange connector.ExchangeName, shares, minPrice numerical.Decimal, expiration int64) PredictionSignalBuilder
	Build() *PredictionSignal
}

// SignalBuilder is the legacy spot signal builder interface, kept for backward compatibility.
// Prefer SpotSignalBuilder for new code.
type SignalBuilder = SpotSignalBuilder
