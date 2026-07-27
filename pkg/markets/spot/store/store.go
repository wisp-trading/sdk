package store

import (
	baseStore "github.com/wisp-trading/sdk/pkg/markets/base/store"
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
	ext := baseStore.NewPairExtensions(timeProvider)
	return &spotStore{
		MarketStore:             ext.MarketStore,
		OrderBookStoreExtension: ext.OrderBookStoreExtension,
		KlineStoreExtension:     ext.KlineStoreExtension,
		TradesStoreExtension:    ext.TradesStoreExtension,
		PositionsStoreExtension: ext.PositionsStoreExtension,
	}
}

func (s *spotStore) MarketType() connector.MarketType {
	return connector.MarketTypeSpot
}
