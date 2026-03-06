package position

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func (ds *dataStore) AddTrade(trade connector.Trade) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.trades.Store(append(ds.getTrades(), trade))
}

func (ds *dataStore) GetTrades() []connector.Trade {
	return ds.getTrades()
}
