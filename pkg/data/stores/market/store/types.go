package store

import (
	"sync"

	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
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

func (ds *dataStore) MarketType() connector.MarketType {
	return connector.MarketTypeUnknown
}
