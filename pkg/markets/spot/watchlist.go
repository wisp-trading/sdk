package spot

import (
	"sync"

	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type pairKey struct {
	Exchange connector.ExchangeName
	Symbol   string
}

type spotWatchlist struct {
	mu       sync.RWMutex
	pairs    map[pairKey]portfolio.Pair
	watchers map[connector.ExchangeName]chan spotTypes.SpotWatchEvent
}

func NewSpotWatchlist() spotTypes.SpotWatchlist {
	return &spotWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan spotTypes.SpotWatchEvent),
	}
}

func (w *spotWatchlist) RequirePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{Exchange: exchange, Symbol: pair.Symbol()}
	if _, exists := w.pairs[key]; exists {
		return
	}
	w.pairs[key] = pair
	w.emitEventLocked(spotTypes.SpotWatchEvent{Exchange: exchange, Pair: pair, Type: spotTypes.SpotPairAdded})
}

func (w *spotWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{Exchange: exchange, Symbol: pair.Symbol()}
	if _, exists := w.pairs[key]; !exists {
		return
	}
	delete(w.pairs, key)
	w.emitEventLocked(spotTypes.SpotWatchEvent{Exchange: exchange, Pair: pair, Type: spotTypes.SpotPairRemoved})
}

func (w *spotWatchlist) GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair {
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

func (w *spotWatchlist) Subscribe(exchange connector.ExchangeName) chan spotTypes.SpotWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}
	ch := make(chan spotTypes.SpotWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *spotWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch, ok := w.watchers[exchange]
	if !ok {
		return
	}
	delete(w.watchers, exchange)
	close(ch)
}

func (w *spotWatchlist) emitEventLocked(ev spotTypes.SpotWatchEvent) {
	ch, ok := w.watchers[ev.Exchange]
	if !ok {
		return
	}
	select {
	case ch <- ev:
	default:
	}
}
