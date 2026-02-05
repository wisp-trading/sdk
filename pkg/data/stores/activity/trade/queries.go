package trade

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// GetAllTrades retrieves all trades
func (ds *dataStore) GetAllTrades() []connector.Trade {
	return ds.getTrades()
}

// GetTradesByExchange retrieves trades for a specific exchange
func (ds *dataStore) GetTradesByExchange(exchange connector.ExchangeName) []connector.Trade {
	trades := ds.getTrades()
	result := make([]connector.Trade, 0)

	for _, trade := range trades {
		if trade.Exchange == exchange {
			result = append(result, trade)
		}
	}

	return result
}

// GetTradesByPair retrieves trades for a specific asset
func (ds *dataStore) GetTradesByPair(asset portfolio.Pair) []connector.Trade {
	trades := ds.getTrades()
	result := make([]connector.Trade, 0)

	for _, trade := range trades {
		if trade.Pair.Symbol() == asset.Symbol() {
			result = append(result, trade)
		}
	}

	return result
}

// GetTradesByExchangeAndPair retrieves trades for a specific exchange and asset
func (ds *dataStore) GetTradesByExchangeAndPair(exchange connector.ExchangeName, pair portfolio.Pair) []connector.Trade {
	trades := ds.getTrades()
	result := make([]connector.Trade, 0)

	for _, trade := range trades {
		if trade.Exchange == exchange && trade.Pair.Symbol() == pair.Symbol() {
			result = append(result, trade)
		}
	}

	return result
}

// GetTradesSince retrieves trades since a specific time
func (ds *dataStore) GetTradesSince(since time.Time) []connector.Trade {
	trades := ds.getTrades()
	result := make([]connector.Trade, 0)

	for _, trade := range trades {
		if trade.Timestamp.After(since) || trade.Timestamp.Equal(since) {
			result = append(result, trade)
		}
	}

	return result
}

// GetTradeByID retrieves a trade by ID
func (ds *dataStore) GetTradeByID(tradeID string) *connector.Trade {
	tradeMap := ds.getTradeMap()

	if trade, exists := tradeMap[tradeID]; exists {
		return &trade
	}

	return nil
}

// TradeExists checks if a trade exists by ID
func (ds *dataStore) TradeExists(tradeID string) bool {
	tradeMap := ds.getTradeMap()
	_, exists := tradeMap[tradeID]
	return exists
}

// GetTradeCount returns the total number of trades
func (ds *dataStore) GetTradeCount() int {
	trades := ds.getTrades()
	return len(trades)
}

// GetTotalVolume calculates total volume for a specific asset
func (ds *dataStore) GetTotalVolume(pair portfolio.Pair) numerical.Decimal {
	trades := ds.getTrades()
	totalVolume := numerical.Zero()

	for _, trade := range trades {
		if trade.Pair.Symbol() == pair.Symbol() {
			totalVolume = totalVolume.Add(trade.Quantity)
		}
	}

	return totalVolume
}
