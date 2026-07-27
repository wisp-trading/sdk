package signal

import (
	baseSignal "github.com/wisp-trading/sdk/pkg/markets/base/signal"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spotBuilder embeds the shared pair builder core.
type spotBuilder struct {
	*baseSignal.Core
}

func (b *spotBuilder) Buy(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.Buy(pair, exchange, quantity)
	return b
}

func (b *spotBuilder) BuyLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.BuyLimit(pair, exchange, quantity, price)
	return b
}

func (b *spotBuilder) Sell(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.Sell(pair, exchange, quantity)
	return b
}

func (b *spotBuilder) SellLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.SellLimit(pair, exchange, quantity, price)
	return b
}

func (b *spotBuilder) SellShort(pair portfolio.Pair, exchange connector.ExchangeName, quantity numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.SellShort(pair, exchange, quantity)
	return b
}

func (b *spotBuilder) SellShortLimit(pair portfolio.Pair, exchange connector.ExchangeName, quantity, price numerical.Decimal) spotTypes.SpotSignalBuilder {
	b.Core.SellShortLimit(pair, exchange, quantity, price)
	return b
}

func (b *spotBuilder) Build() (spotTypes.SpotSignal, error) {
	if err := b.ValidateReady(); err != nil {
		return nil, err
	}
	return spotTypes.NewSpotSignal(b.NewID(), b.StrategyName, b.TimeProvider.Now(), b.Actions), nil
}
