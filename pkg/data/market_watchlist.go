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

	// watchers are keyed by exchange, each with its own set of channels
	watchers map[connector.ExchangeName]map[chan data.MarketWatchEvent]struct{}
}

func NewMarketWatchlist() data.MarketWatchlist {
	return &marketWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]map[chan data.MarketWatchEvent]struct{}),
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

func (w *marketWatchlist) Subscribe(exchange connector.ExchangeName) <-chan data.MarketWatchEvent {
	ch := make(chan data.MarketWatchEvent, 128)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.watchers[exchange] == nil {
		w.watchers[exchange] = make(map[chan data.MarketWatchEvent]struct{})
	}
	w.watchers[exchange][ch] = struct{}{}

	return ch
}

func (w *marketWatchlist) Unsubscribe(exchange connector.ExchangeName, ch chan data.MarketWatchEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	watchersForEx, ok := w.watchers[exchange]
	if !ok {
		return
	}

	if _, ok := watchersForEx[ch]; ok {
		delete(watchersForEx, ch)
		close(ch)
	}

	if len(watchersForEx) == 0 {
		delete(w.watchers, exchange)
	}
}

func (w *marketWatchlist) emitEventLocked(ev data.MarketWatchEvent) {
	ex := ev.Requirement.Exchange

	watchersForEx, ok := w.watchers[ex]
	if !ok {
		return
	}

	for ch := range watchersForEx {
		select {
		case ch <- ev:
		default:
			// drop or log slow watcher
		}
	}
}
