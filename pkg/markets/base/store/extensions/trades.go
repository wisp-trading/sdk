package extensions

import (
	"sync"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type tradesExtension struct {
	mu     sync.RWMutex
	trades []connector.Trade
	byID   map[string]int // tradeID -> index
}

func NewTradesExtension() market.TradesStoreExtension {
	return &tradesExtension{
		trades: []connector.Trade{},
		byID:   make(map[string]int),
	}
}

func (e *tradesExtension) AddTrade(trade connector.Trade) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.byID[trade.ID]; exists {
		return
	}
	e.byID[trade.ID] = len(e.trades)
	e.trades = append(e.trades, trade)
}

func (e *tradesExtension) AddTrades(trades []connector.Trade) {
	for _, t := range trades {
		e.AddTrade(t)
	}
}

func (e *tradesExtension) GetAllTrades() []connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]connector.Trade, len(e.trades))
	copy(out, e.trades)
	return out
}

func (e *tradesExtension) GetTradesByExchange(exchange connector.ExchangeName) []connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []connector.Trade
	for _, t := range e.trades {
		if t.Exchange == exchange {
			out = append(out, t)
		}
	}
	return out
}

func (e *tradesExtension) GetTradesByPair(pair portfolio.Pair) []connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []connector.Trade
	for _, t := range e.trades {
		if t.Pair.Symbol() == pair.Symbol() {
			out = append(out, t)
		}
	}
	return out
}

func (e *tradesExtension) GetTradesSince(since time.Time) []connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []connector.Trade
	for _, t := range e.trades {
		if !t.Timestamp.Before(since) {
			out = append(out, t)
		}
	}
	return out
}

func (e *tradesExtension) GetTradeByID(tradeID string) *connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if idx, ok := e.byID[tradeID]; ok {
		t := e.trades[idx]
		return &t
	}
	return nil
}

func (e *tradesExtension) TradeExists(tradeID string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.byID[tradeID]
	return ok
}

func (e *tradesExtension) GetTradeCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.trades)
}

func (e *tradesExtension) GetTotalVolume(pair portfolio.Pair) numerical.Decimal {
	e.mu.RLock()
	defer e.mu.RUnlock()
	total := numerical.Zero()
	for _, t := range e.trades {
		if t.Pair.Symbol() == pair.Symbol() {
			total = total.Add(t.Quantity)
		}
	}
	return total
}

func (e *tradesExtension) QueryTrades(q market.ActivityQuery) []connector.Trade {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []connector.Trade
	for _, t := range e.trades {
		if q.Exchange != nil && t.Exchange != *q.Exchange {
			continue
		}
		if q.Pair != nil && t.Pair.Symbol() != q.Pair.Symbol() {
			continue
		}
		out = append(out, t)
	}
	return out
}

var _ market.TradesStoreExtension = (*tradesExtension)(nil)
