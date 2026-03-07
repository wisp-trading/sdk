package store

import (
	baseStore "github.com/wisp-trading/sdk/pkg/markets/base/store"
	"github.com/wisp-trading/sdk/pkg/markets/base/store/extensions"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type spotStore struct {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	market.PositionsStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) spotTypes.MarketStore {
	base := baseStore.NewStore(timeProvider)

	return &spotStore{
		MarketStore:             base,
		OrderBookStoreExtension: extensions.NewOrderBookExtension(base.UpdatePairPrice, base.UpdateLastUpdated),
		KlineStoreExtension:     extensions.NewKlineExtension(),
		TradesStoreExtension:    extensions.NewTradesExtension(),
		PositionsStoreExtension: extensions.NewPositionsExtension(),
	}
}

func (s *spotStore) MarketType() connector.MarketType {
	return connector.MarketTypeSpot
}
