package perp

import (
	"sync"

	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type pairKey struct {
	Exchange connector.ExchangeName
	Symbol   string
}

type perpWatchlist struct {
	mu       sync.RWMutex
	pairs    map[pairKey]portfolio.Pair
	watchers map[connector.ExchangeName]chan baseTypes.MarketWatchEvent
}

// NewPerpWatchlist creates a new empty perp-domain watchlist.
func NewPerpWatchlist() perpTypes.PerpWatchlist {
	return &perpWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan baseTypes.MarketWatchEvent),
	}
}

func (w *perpWatchlist) RequirePair(exchange connector.ExchangeName, pair portfolio.Pair) {
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

func (w *perpWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
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

func (w *perpWatchlist) GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair {
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

func (w *perpWatchlist) Subscribe(exchange connector.ExchangeName) chan baseTypes.MarketWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}
	ch := make(chan baseTypes.MarketWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *perpWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		delete(w.watchers, exchange)
		close(ch)
	}
}

func (w *perpWatchlist) emitLocked(ev baseTypes.MarketWatchEvent) {
	if ch, ok := w.watchers[ev.Requirement.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

var _ perpTypes.PerpWatchlist = (*perpWatchlist)(nil)
