package store

import (
	baseStore "github.com/wisp-trading/sdk/pkg/markets/base/store"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	perpExtensions "github.com/wisp-trading/sdk/pkg/markets/perp/store/extensions"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type perpStore struct {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
	market.TradesStoreExtension
	market.PositionsStoreExtension
	domainTypes.FundingRateStoreExtension
	domainTypes.PerpPositionsStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) domainTypes.MarketStore {
	ext := baseStore.NewPairExtensions(timeProvider)
	return &perpStore{
		MarketStore:                 ext.MarketStore,
		OrderBookStoreExtension:     ext.OrderBookStoreExtension,
		KlineStoreExtension:         ext.KlineStoreExtension,
		TradesStoreExtension:        ext.TradesStoreExtension,
		PositionsStoreExtension:     ext.PositionsStoreExtension,
		FundingRateStoreExtension:   perpExtensions.NewFundingRateExtension(),
		PerpPositionsStoreExtension: perpExtensions.NewPerpPositionsExtension(),
	}
}

func (ps *perpStore) MarketType() connector.MarketType {
	return connector.MarketTypePerp
}
