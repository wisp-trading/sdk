// Package views provides shared pair-market monitoring helpers for spot/perp.
package views

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PairRef is one watched exchange+pair (mapped to SpotMarketView / PerpMarketView by domains).
type PairRef struct {
	Exchange connector.ExchangeName
	Pair     portfolio.Pair
}

// BookReader is the store surface for orderbook/klines used by monitoring views.
type BookReader interface {
	GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetKlines(pair portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
}

// ListWatchedPairs returns every pair on the watchlist for the given ready exchanges.
func ListWatchedPairs(
	readyNames []connector.ExchangeName,
	watchlist baseTypes.MarketWatchlist,
) []PairRef {
	result := make([]PairRef, 0)
	for _, name := range readyNames {
		for _, pair := range watchlist.GetRequiredPairs(name) {
			result = append(result, PairRef{Exchange: name, Pair: pair})
		}
	}
	return result
}

// OrderBook reads from a pair store (nil if missing).
func OrderBook(store BookReader, exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook {
	if store == nil {
		return nil
	}
	return store.GetOrderBook(pair, exchange)
}

// Klines reads from a pair store.
func Klines(store BookReader, exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	if store == nil {
		return nil
	}
	return store.GetKlines(pair, exchange, interval, limit)
}
