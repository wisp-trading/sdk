package data

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// internal key for the map.
type pairKey struct {
	Exchange connector.ExchangeName
	Symbol   string
}

// marketWatchlist is the concrete implementation of MarketWatchlist.
type marketWatchlist struct {
	mu sync.RWMutex

	pairs map[pairKey]portfolio.Pair

	// exactly one channel per exchange
	watchers map[connector.ExchangeName]chan data.MarketWatchEvent
}

func NewMarketWatchlist() data.MarketWatchlist {
	return &marketWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan data.MarketWatchEvent),
	}
}

func (w *marketWatchlist) RequirePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{
		Exchange: exchange,
		Symbol:   pair.Symbol(),
	}

	if _, exists := w.pairs[key]; exists {
		// Already required; no-op.
		return
	}

	w.pairs[key] = pair

	w.emitEventLocked(data.MarketWatchEvent{
		Requirement: data.PairRequirement{
			Exchange: exchange,
			Pair:     pair,
		},
		Type: data.PairAdded,
	})
}

func (w *marketWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := pairKey{
		Exchange: exchange,
		Symbol:   pair.Symbol(),
	}

	if _, exists := w.pairs[key]; !exists {
		// Not present; no-op.
		return
	}

	delete(w.pairs, key)

	w.emitEventLocked(data.MarketWatchEvent{
		Requirement: data.PairRequirement{
			Exchange: exchange,
			Pair:     pair,
		},
		Type: data.PairRemoved,
	})
}

func (w *marketWatchlist) GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair {
	w.mu.RLock()
	defer w.mu.RUnlock()

	out := make([]portfolio.Pair, 0, len(w.pairs))
	for key, pair := range w.pairs {
		if key.Exchange == exchange {
			out = append(out, pair)
		}
	}

	return out
}

func (w *marketWatchlist) Subscribe(exchange connector.ExchangeName) chan data.MarketWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	// If a channel already exists for this exchange, reuse it.
	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}

	ch := make(chan data.MarketWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *marketWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch, ok := w.watchers[exchange]
	if !ok {
		return
	}

	delete(w.watchers, exchange)
	close(ch)
}

func (w *marketWatchlist) emitEventLocked(ev data.MarketWatchEvent) {
	ex := ev.Requirement.Exchange

	ch, ok := w.watchers[ex]
	if !ok {
		return
	}

	select {
	case ch <- ev:
	default:
		// drop or log slow watcher
	}
}
