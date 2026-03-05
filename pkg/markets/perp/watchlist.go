package perp

import (
	"sync"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// internal key for the map.
type pairKey struct {
	Exchange connector.ExchangeName
	Symbol   string
}

// perpWatchlist is the concrete implementation of PerpWatchlist.
type perpWatchlist struct {
	mu sync.RWMutex

	pairs map[pairKey]portfolio.Pair

	// exactly one channel per exchange
	watchers map[connector.ExchangeName]chan perpTypes.PerpWatchEvent
}

// NewPerpWatchlist creates a new perp-domain watchlist.
func NewPerpWatchlist() perpTypes.PerpWatchlist {
	return &perpWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan perpTypes.PerpWatchEvent),
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

	w.emitEventLocked(perpTypes.PerpWatchEvent{
		Exchange: exchange,
		Pair:     pair,
		Type:     perpTypes.PerpPairAdded,
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

	w.emitEventLocked(perpTypes.PerpWatchEvent{
		Exchange: exchange,
		Pair:     pair,
		Type:     perpTypes.PerpPairRemoved,
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

func (w *perpWatchlist) Subscribe(exchange connector.ExchangeName) chan perpTypes.PerpWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}

	ch := make(chan perpTypes.PerpWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *perpWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch, ok := w.watchers[exchange]
	if !ok {
		return
	}

	delete(w.watchers, exchange)
	close(ch)
}

func (w *perpWatchlist) emitEventLocked(ev perpTypes.PerpWatchEvent) {
	ch, ok := w.watchers[ev.Exchange]
	if !ok {
		return
	}

	select {
	case ch <- ev:
	default:
		// drop on slow watcher
	}
}
