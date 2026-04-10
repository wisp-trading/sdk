package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// OptionsSignalBuilder provides a fluent API for constructing options market signals.
// Obtain one via wisp.Options().Signal(strategyName).
type OptionsSignalBuilder interface {
	// Buy opens a long position on an options contract (buy to open).
	Buy(contract OptionContract, exchange connector.ExchangeName, quantity numerical.Decimal) OptionsSignalBuilder

	// BuyLimit opens a long position at a specified limit price.
	BuyLimit(contract OptionContract, exchange connector.ExchangeName, quantity, price numerical.Decimal) OptionsSignalBuilder

	// Sell opens a short position or closes a long position (sell to open / sell to close).
	Sell(contract OptionContract, exchange connector.ExchangeName, quantity numerical.Decimal) OptionsSignalBuilder

	// SellLimit opens a short or closes a long at a specified limit price.
	SellLimit(contract OptionContract, exchange connector.ExchangeName, quantity, price numerical.Decimal) OptionsSignalBuilder

	Build() (OptionsSignal, error)
}
