package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// OnchainSignalBuilder builds exact-in swap signals.
//
// Buy(quoteAmount) — spend quote for base.
// Sell(baseAmount) — sell base for quote.
// Limit methods exist for API symmetry but the UniV3 pilot rejects non-zero limit prices.
type OnchainSignalBuilder interface {
	Buy(pair portfolio.Pair, exchange connector.ExchangeName, quoteAmount numerical.Decimal) OnchainSignalBuilder
	BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) OnchainSignalBuilder
	Sell(pair portfolio.Pair, exchange connector.ExchangeName, baseAmount numerical.Decimal) OnchainSignalBuilder
	SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) OnchainSignalBuilder
	Build() (OnchainSignal, error)
}
