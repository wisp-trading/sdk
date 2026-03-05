package store

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	storeExtensions "github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store/extensions"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type perpStore struct {
	marketTypes.MarketStore
	marketTypes.OrderBookStoreExtension
	marketTypes.KlineStoreExtension
	domainTypes.FundingRateStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) domainTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &perpStore{
		MarketStore: baseStore,
		OrderBookStoreExtension: storeExtensions.NewOrderBookExtension(
			baseStore.UpdatePairPrice,
			baseStore.UpdateLastUpdated,
		),
		KlineStoreExtension:       storeExtensions.NewKlineExtension(),
		FundingRateStoreExtension: extensions.NewFundingRateExtension(),
	}
}

func (ps *perpStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypePerp
}
