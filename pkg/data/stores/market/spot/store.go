package spot

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	storeExtensions "github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	spotTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// spotStore composes base store with extensions
type spotStore struct {
	marketTypes.MarketStore
	marketTypes.OrderBookStoreExtension
	marketTypes.KlineStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) spotTypes.MarketStore {
	// Create base store first
	baseStore := store.NewStore(timeProvider)

	return &spotStore{
		MarketStore: baseStore,
		OrderBookStoreExtension: storeExtensions.NewOrderBookExtension(
			baseStore.UpdatePairPrice,
			baseStore.UpdateLastUpdated,
		),
		KlineStoreExtension: storeExtensions.NewKlineExtension(
			baseStore.UpdateLastUpdated,
		),
	}
}

func (s *spotStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypeSpot
}
