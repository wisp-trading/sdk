package trade

import "github.com/wisp-trading/sdk/pkg/types/connector"

func (ds *dataStore) AddTrade(trade connector.Trade) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	tradeMap := ds.getTradeMap()
	if _, exists := tradeMap[trade.ID]; exists {
		return
	}

	currentTrades := ds.getTrades()
	updatedTrades := make([]connector.Trade, len(currentTrades)+1)
	copy(updatedTrades, currentTrades)
	updatedTrades[len(currentTrades)] = trade
	ds.trades.Store(updatedTrades)

	updatedMap := make(TradeMap, len(tradeMap)+1)
	for k, v := range tradeMap {
		updatedMap[k] = v
	}
	updatedMap[trade.ID] = trade
	ds.byID.Store(updatedMap)
}

func (ds *dataStore) AddTrades(trades []connector.Trade) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	tradeMap := ds.getTradeMap()
	currentTrades := ds.getTrades()

	newTrades := make([]connector.Trade, 0)
	for _, trade := range trades {
		if _, exists := tradeMap[trade.ID]; !exists {
			newTrades = append(newTrades, trade)
		}
	}
	if len(newTrades) == 0 {
		return
	}

	updatedTrades := make([]connector.Trade, len(currentTrades)+len(newTrades))
	copy(updatedTrades, currentTrades)
	copy(updatedTrades[len(currentTrades):], newTrades)
	ds.trades.Store(updatedTrades)

	updatedMap := make(TradeMap, len(tradeMap)+len(newTrades))
	for k, v := range tradeMap {
		updatedMap[k] = v
	}
	for _, trade := range newTrades {
		updatedMap[trade.ID] = trade
	}
	ds.byID.Store(updatedMap)
}
