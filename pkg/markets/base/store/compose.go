package store

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/store/extensions"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// PairExtensions is the shared store surface for pair domains (spot/perp).
type PairExtensions struct {
	MarketStore             market.MarketStore
	OrderBookStoreExtension market.OrderBookStoreExtension
	KlineStoreExtension     market.KlineStoreExtension
	TradesStoreExtension    market.TradesStoreExtension
	PositionsStoreExtension market.PositionsStoreExtension
}

// NewPairExtensions builds base market store + orderbook/kline/trades/positions extensions.
// Domain stores embed this and add market-specific extensions (funding, etc.).
func NewPairExtensions(timeProvider temporal.TimeProvider) PairExtensions {
	base := NewStore(timeProvider)
	return PairExtensions{
		MarketStore:             base,
		OrderBookStoreExtension: extensions.NewOrderBookExtension(base.UpdatePairPrice, base.UpdateLastUpdated),
		KlineStoreExtension:     extensions.NewKlineExtension(),
		TradesStoreExtension:    extensions.NewTradesExtension(),
		PositionsStoreExtension: extensions.NewPositionsExtension(),
	}
}
