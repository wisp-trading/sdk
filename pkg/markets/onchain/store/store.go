package store

import (
	baseStore "github.com/wisp-trading/sdk/pkg/markets/base/store"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type onchainStore struct {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	market.PositionsStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) onchainTypes.MarketStore {
	ext := baseStore.NewPairExtensions(timeProvider)
	return &onchainStore{
		MarketStore:             ext.MarketStore,
		OrderBookStoreExtension: ext.OrderBookStoreExtension,
		KlineStoreExtension:     ext.KlineStoreExtension,
		TradesStoreExtension:    ext.TradesStoreExtension,
		PositionsStoreExtension: ext.PositionsStoreExtension,
	}
}

func (s *onchainStore) MarketType() connector.MarketType {
	return connector.MarketTypeOnchain
}
