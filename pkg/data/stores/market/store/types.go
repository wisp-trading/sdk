package store

import (
	"sync"
	"sync/atomic"

	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// Data types for atomic storage
type assetOrderBooks map[portfolio.Pair]marketTypes.OrderBookMap
type assetPrices map[portfolio.Pair]marketTypes.PriceMap
type assetKlines map[portfolio.Pair]marketTypes.KlineMap

type dataStore struct {
	timeProvider temporal.TimeProvider
	orderBooks   atomic.Value // assetOrderBooks
	prices       atomic.Value // assetPrices
	klines       atomic.Value // assetKlines
	lastUpdated  atomic.Value // marketTypes.LastUpdatedMap

	mutex                sync.RWMutex
	orchestratorNotifier func()

	// Extensions for market-specific data (funding rates, etc.)
	extensions []marketTypes.StoreExtension
}

func (ds *dataStore) MarketType() marketTypes.MarketType {
	return marketTypes.MarketTypeUnknown
}
