package spot

import (
	"github.com/wisp-trading/wisp/pkg/data/stores/market/store"
	marketTypes "github.com/wisp-trading/wisp/pkg/types/data/stores/market"
	spotTypes "github.com/wisp-trading/wisp/pkg/types/data/stores/market/spot"
	"github.com/wisp-trading/wisp/pkg/types/temporal"
)

type spotStore struct {
	marketTypes.MarketStore
}

func NewStore(timeProvider temporal.TimeProvider) spotTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &spotStore{
		MarketStore: baseStore,
	}
}

func (s *spotStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypeSpot
}
