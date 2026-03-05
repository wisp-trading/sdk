package store

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store/extensions"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type predictionStore struct {
	marketTypes.MarketStore
	domainTypes.OrderBookStoreExtension
	domainTypes.PositionsStoreExtension
	domainTypes.BalanceStoreExtension
}

func NewStore(timeProvider temporal.TimeProvider) domainTypes.MarketStore {
	baseStore := store.NewStore(timeProvider)

	return &predictionStore{
		MarketStore:             baseStore,
		OrderBookStoreExtension: extensions.NewPredictionOrderBookExtension(),
		PositionsStoreExtension: extensions.NewPredictionPositionsExtension(),
		BalanceStoreExtension:   extensions.NewPredictionBalanceExtension(),
	}
}

func (s *predictionStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypePrediction
}
