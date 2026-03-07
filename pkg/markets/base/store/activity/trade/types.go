package trade

import (
	"sync"
	"sync/atomic"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

type TradeMap map[string]connector.Trade

type dataStore struct {
	mutex  sync.RWMutex
	trades atomic.Value // []connector.Trade
	byID   atomic.Value // TradeMap
}

func (ds *dataStore) getTrades() []connector.Trade {
	if v := ds.trades.Load(); v != nil {
		return v.([]connector.Trade)
	}
	return []connector.Trade{}
}

func (ds *dataStore) getTradeMap() TradeMap {
	if v := ds.byID.Load(); v != nil {
		return v.(TradeMap)
	}
	return make(TradeMap)
}
