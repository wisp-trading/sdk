package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// perpBuilder embeds the shared pair builder core; leverage methods are extensions.
type perpBuilder struct {
	*baseSignal.Core
}

func (b *perpBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.Buy(pair, exchange, quantity)
	return b
}

func (b *perpBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.BuyLimit(pair, exchange, quantity, price)
	return b
}

func (b *perpBuilder) BuyLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.BuyLimitWithLeverage(pair, exchange, quantity, price, leverage)
	return b
}

func (b *perpBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.Sell(pair, exchange, quantity)
	return b
}

func (b *perpBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.SellLimit(pair, exchange, quantity, price)
	return b
}

func (b *perpBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.SellShort(pair, exchange, quantity)
	return b
}

func (b *perpBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.SellShortLimit(pair, exchange, quantity, price)
	return b
}

func (b *perpBuilder) SellShortLimitWithLeverage(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price, leverage numerical.Decimal) perpTypes.PerpSignalBuilder {
	b.Core.SellShortLimitWithLeverage(pair, exchange, quantity, price, leverage)
	return b
}

func (b *perpBuilder) Build() (perpTypes.PerpSignal, error) {
	if err := b.ValidateReady(); err != nil {
		return nil, err
	}
	return perpTypes.NewPerpSignal(b.NewID(), b.StrategyName, b.TimeProvider.Now(), b.Actions), nil
}
