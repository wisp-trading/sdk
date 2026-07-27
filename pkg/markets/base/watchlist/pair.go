// Package watchlist provides the shared pair-based MarketWatchlist used by
// spot and perp domains (prediction/options use domain-specific keys).
package watchlist

import (
	"sync"

	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type pairKey struct {
	Exchange connector.ExchangeName
	Symbol   string
}

type pairWatchlist struct {
	mu       sync.RWMutex
	pairs    map[pairKey]portfolio.Pair
	watchers map[connector.ExchangeName]chan baseTypes.MarketWatchEvent
}

// NewPairWatchlist creates an empty pair watchlist (spot/perp domains).
func NewPairWatchlist() baseTypes.MarketWatchlist {
	return &pairWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan baseTypes.MarketWatchEvent),
	}
}

func (w *pairWatchlist) RequirePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{Exchange: exchange, Symbol: pair.Symbol()}
	if _, exists := w.pairs[key]; exists {
		return
	}
	w.pairs[key] = pair
	w.emitLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairAdded,
	})
}

func (w *pairWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{Exchange: exchange, Symbol: pair.Symbol()}
	if _, exists := w.pairs[key]; !exists {
		return
	}
	delete(w.pairs, key)
	w.emitLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairRemoved,
	})
}

func (w *pairWatchlist) GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair {
	w.mu.RLock()
	defer w.mu.RUnlock()

	out := make([]portfolio.Pair, 0)
	for key, pair := range w.pairs {
		if key.Exchange == exchange {
			out = append(out, pair)
		}
	}
	return out
}

func (w *pairWatchlist) Subscribe(exchange connector.ExchangeName) chan baseTypes.MarketWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}
	ch := make(chan baseTypes.MarketWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *pairWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		delete(w.watchers, exchange)
		close(ch)
	}
}

func (w *pairWatchlist) emitLocked(ev baseTypes.MarketWatchEvent) {
	if ch, ok := w.watchers[ev.Requirement.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

var _ baseTypes.MarketWatchlist = (*pairWatchlist)(nil)
