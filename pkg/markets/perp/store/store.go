package store

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	storeExtensions "github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/perp/store/extensions"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type perpStore struct {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
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

func (ps *perpStore) MarketType() connector.MarketType {
	return connector.MarketTypePerp
}
