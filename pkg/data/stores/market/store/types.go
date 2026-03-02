package store

import (
	"sync"

	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type dataStore struct {
	timeProvider temporal.TimeProvider
	prices       map[portfolio.Pair]marketTypes.PriceMap
	lastUpdated  marketTypes.LastUpdatedMap

	mu                   sync.RWMutex
	orchestratorNotifier func()

	extensions []marketTypes.StoreExtension
}

func (ds *dataStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypeUnknown
}
