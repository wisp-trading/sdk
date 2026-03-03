package store

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store/extensions"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	predictionTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type predictionStore struct {
	marketTypes.MarketStore
	predictionTypes.OrderBookStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) predictionTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &predictionStore{
		MarketStore:             baseStore,
		OrderBookStoreExtension: extensions.NewPredictionOrderBookExtension(),
	}
}

func (s *predictionStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypePrediction
}
