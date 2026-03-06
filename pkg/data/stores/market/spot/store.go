package spot

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	storeExtensions "github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// spotStore composes base store with extensions
type spotStore struct {
	market.MarketStore
	market.OrderBookStoreExtension
	market.KlineStoreExtension
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
		KlineStoreExtension: storeExtensions.NewKlineExtension(),
	}
}

func (s *spotStore) MarketType() connector.MarketType {
	return connector.MarketTypeSpot
}
