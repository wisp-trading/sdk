package store

import (
	"sync"
	"sync/atomic"

	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// Data types for atomic storage
type assetPrices map[portfolio.Pair]marketTypes.PriceMap

type dataStore struct {
	timeProvider temporal.TimeProvider
	prices       atomic.Value // assetPrices
	lastUpdated  atomic.Value // marketTypes.LastUpdatedMap

	mutex                sync.RWMutex
	orchestratorNotifier func()

	extensions []marketTypes.StoreExtension
}

func (ds *dataStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypeUnknown
}
