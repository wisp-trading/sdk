package trade

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func NewStore() activity.Trades {
	ds := &dataStore{}
	ds.trades.Store([]connector.Trade{})
	ds.byID.Store(make(TradeMap))
	return ds
}

// Clear removes all trades from the store (for simulation restart)
func (ds *dataStore) Clear() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	ds.trades.Store([]connector.Trade{})
	ds.byID.Store(make(TradeMap))
}

var _ activity.Trades = (*dataStore)(nil)
