package store

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/store"
	"github.com/wisp-trading/sdk/pkg/markets/base/store/extensions"
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
	domainTypes.FundingRateStoreExtension
	domainTypes.PerpPositionsStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) domainTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &perpStore{
		MarketStore:                 baseStore,
		OrderBookStoreExtension:     extensions.NewOrderBookExtension(baseStore.UpdatePairPrice, baseStore.UpdateLastUpdated),
		KlineStoreExtension:         extensions.NewKlineExtension(),
		TradesStoreExtension:        extensions.NewTradesExtension(),
		FundingRateStoreExtension:   perpExtensions.NewFundingRateExtension(),
		PerpPositionsStoreExtension: perpExtensions.NewPerpPositionsExtension(),
	}
}

func (ps *perpStore) MarketType() connector.MarketType {
	return connector.MarketTypePerp
}
