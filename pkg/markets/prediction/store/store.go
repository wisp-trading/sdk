package store

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/store"
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/store/extensions"
	domainTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
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

func (s *predictionStore) MarketType() connector.MarketType {
	return connector.MarketTypePrediction
}
