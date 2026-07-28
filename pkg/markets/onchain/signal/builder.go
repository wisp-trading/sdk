package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type onchainBuilder struct {
	*baseSignal.Core
}

func (b *onchainBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quoteAmount numerical.Decimal) onchainTypes.OnchainSignalBuilder {
	b.Core.Buy(pair, exchange, quoteAmount)
	return b
}

func (b *onchainBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) onchainTypes.OnchainSignalBuilder {
	b.Core.BuyLimit(pair, exchange, quantity, price)
	return b
}

func (b *onchainBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, baseAmount numerical.Decimal) onchainTypes.OnchainSignalBuilder {
	b.Core.Sell(pair, exchange, baseAmount)
	return b
}

func (b *onchainBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) onchainTypes.OnchainSignalBuilder {
	b.Core.SellLimit(pair, exchange, quantity, price)
	return b
}

func (b *onchainBuilder) Build() (onchainTypes.OnchainSignal, error) {
	if err := b.ValidateReady(); err != nil {
		return nil, err
	}
	return onchainTypes.NewOnchainSignal(b.NewID(), b.StrategyName, b.TimeProvider.Now(), b.Actions), nil
}

var _ onchainTypes.OnchainSignalBuilder = (*onchainBuilder)(nil)
