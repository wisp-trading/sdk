package perp

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/perp/extensions"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	storeExtensions "github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type perpStore struct {
	marketTypes.MarketStore
	marketTypes.OrderBookStoreExtension
	marketTypes.KlineStoreExtension
	perpTypes.FundingRateStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) perpTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &perpStore{
		MarketStore: baseStore,
		OrderBookStoreExtension: storeExtensions.NewOrderBookExtension(
			baseStore.UpdatePairPrice,
			baseStore.UpdateLastUpdated,
		),
		KlineStoreExtension: storeExtensions.NewKlineExtension(
			baseStore.UpdateLastUpdated,
		),
		FundingRateStoreExtension: extensions.NewFundingRateExtension(),
	}
}

func (ps *perpStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypePerp
}
