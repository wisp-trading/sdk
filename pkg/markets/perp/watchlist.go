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
	mu           sync.RWMutex
	pairs        map[pairKey]portfolio.Pair
	baseWatchers map[connector.ExchangeName]chan baseTypes.MarketWatchEvent
	perpWatchers map[connector.ExchangeName]chan perpTypes.PerpWatchEvent
}

// NewPerpWatchlist creates a new perp-domain watchlist.
func NewPerpWatchlist() perpTypes.PerpWatchlist {
	return &perpWatchlist{
		pairs:        make(map[pairKey]portfolio.Pair),
		baseWatchers: make(map[connector.ExchangeName]chan baseTypes.MarketWatchEvent),
		perpWatchers: make(map[connector.ExchangeName]chan perpTypes.PerpWatchEvent),
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
	w.emitBaseLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairAdded,
	})
	w.emitPerpLocked(perpTypes.PerpWatchEvent{Exchange: exchange, Pair: pair, Type: perpTypes.PerpPairAdded})
}

func (w *perpWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := pairKey{Exchange: exchange, Symbol: pair.Symbol()}
	if _, exists := w.pairs[key]; !exists {
		return
	}
	delete(w.pairs, key)
	w.emitBaseLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairRemoved,
	})
	w.emitPerpLocked(perpTypes.PerpWatchEvent{Exchange: exchange, Pair: pair, Type: perpTypes.PerpPairRemoved})
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
	if ch, ok := w.baseWatchers[exchange]; ok {
		return ch
	}
	ch := make(chan baseTypes.MarketWatchEvent, 128)
	w.baseWatchers[exchange] = ch
	return ch
}

func (w *perpWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if ch, ok := w.baseWatchers[exchange]; ok {
		delete(w.baseWatchers, exchange)
		close(ch)
	}
}

func (w *perpWatchlist) SubscribePerp(exchange connector.ExchangeName) chan perpTypes.PerpWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()
	if ch, ok := w.perpWatchers[exchange]; ok {
		return ch
	}
	ch := make(chan perpTypes.PerpWatchEvent, 128)
	w.perpWatchers[exchange] = ch
	return ch
}

func (w *perpWatchlist) UnsubscribePerp(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if ch, ok := w.perpWatchers[exchange]; ok {
		delete(w.perpWatchers, exchange)
		close(ch)
	}
}

func (w *perpWatchlist) emitBaseLocked(ev baseTypes.MarketWatchEvent) {
	if ch, ok := w.baseWatchers[ev.Requirement.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

func (w *perpWatchlist) emitPerpLocked(ev perpTypes.PerpWatchEvent) {
	if ch, ok := w.perpWatchers[ev.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

var _ perpTypes.PerpWatchlist = (*perpWatchlist)(nil)
