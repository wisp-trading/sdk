package market

import (
	"sync"
	"sync/atomic"

	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// Data types for atomic storage
type assetFundingRates map[portfolio.Asset]marketTypes.FundingRateMap
type assetHistoricalFunding map[portfolio.Asset]marketTypes.HistoricalFundingMap
type assetOrderBooks map[portfolio.Asset]marketTypes.OrderBookMap
type assetPrices map[portfolio.Asset]marketTypes.PriceMap
type assetKlines map[portfolio.Asset]marketTypes.KlineMap

type dataStore struct {
	timeProvider           temporal.TimeProvider
	fundingRates           atomic.Value // assetFundingRates
	historicalFundingRates atomic.Value // assetHistoricalFunding
	orderBooks             atomic.Value // assetOrderBooks
	prices                 atomic.Value // assetPrices
	klines                 atomic.Value // assetKlines
	lastUpdated            atomic.Value // marketTypes.LastUpdatedMap

	mutex                sync.RWMutex
	orchestratorNotifier func()
}
