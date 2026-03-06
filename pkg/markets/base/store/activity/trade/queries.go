package trade

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

func (ds *dataStore) GetAllTrades() []connector.Trade {
	return ds.getTrades()
}

func (ds *dataStore) GetTradesByExchange(exchange connector.ExchangeName) []connector.Trade {
	result := make([]connector.Trade, 0)
	for _, trade := range ds.getTrades() {
		if trade.Exchange == exchange {
			result = append(result, trade)
		}
	}
	return result
}

func (ds *dataStore) GetTradesByPair(asset portfolio.Pair) []connector.Trade {
	result := make([]connector.Trade, 0)
	for _, trade := range ds.getTrades() {
		if trade.Pair.Symbol() == asset.Symbol() {
			result = append(result, trade)
		}
	}
	return result
}

func (ds *dataStore) GetTradesByExchangeAndPair(exchange connector.ExchangeName, pair portfolio.Pair) []connector.Trade {
	result := make([]connector.Trade, 0)
	for _, trade := range ds.getTrades() {
		if trade.Exchange == exchange && trade.Pair.Symbol() == pair.Symbol() {
			result = append(result, trade)
		}
	}
	return result
}

func (ds *dataStore) GetTradesSince(since time.Time) []connector.Trade {
	result := make([]connector.Trade, 0)
	for _, trade := range ds.getTrades() {
		if trade.Timestamp.After(since) || trade.Timestamp.Equal(since) {
			result = append(result, trade)
		}
	}
	return result
}

func (ds *dataStore) GetTradeByID(tradeID string) *connector.Trade {
	if trade, exists := ds.getTradeMap()[tradeID]; exists {
		return &trade
	}
	return nil
}

func (ds *dataStore) TradeExists(tradeID string) bool {
	_, exists := ds.getTradeMap()[tradeID]
	return exists
}

func (ds *dataStore) GetTradeCount() int {
	return len(ds.getTrades())
}

func (ds *dataStore) GetTotalVolume(pair portfolio.Pair) numerical.Decimal {
	total := numerical.Zero()
	for _, trade := range ds.getTrades() {
		if trade.Pair.Symbol() == pair.Symbol() {
			total = total.Add(trade.Quantity)
		}
	}
	return total
}
