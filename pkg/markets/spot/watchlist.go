package spot

import (
	"sync"

	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
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
	watchers map[connector.ExchangeName]chan baseTypes.MarketWatchEvent
}

// NewSpotWatchlist creates a new empty spot-domain watchlist.
// In the fx graph the module seeds it with config assets; in tests it can be
// used directly and seeded via RequirePair.
func NewSpotWatchlist() spotTypes.SpotWatchlist {
	return &spotWatchlist{
		pairs:    make(map[pairKey]portfolio.Pair),
		watchers: make(map[connector.ExchangeName]chan baseTypes.MarketWatchEvent),
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
	w.emitLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairAdded,
	})
}

func (w *spotWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
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

// Subscribe satisfies base MarketWatchlist — used by the base ingestor.
func (w *spotWatchlist) Subscribe(exchange connector.ExchangeName) chan baseTypes.MarketWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}
	ch := make(chan baseTypes.MarketWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *spotWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		delete(w.watchers, exchange)
		close(ch)
	}
}

func (w *spotWatchlist) emitLocked(ev baseTypes.MarketWatchEvent) {
	if ch, ok := w.watchers[ev.Requirement.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

var _ spotTypes.SpotWatchlist = (*spotWatchlist)(nil)
