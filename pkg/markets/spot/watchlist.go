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
	mu           sync.RWMutex
	pairs        map[pairKey]portfolio.Pair
	baseWatchers map[connector.ExchangeName]chan baseTypes.MarketWatchEvent
	spotWatchers map[connector.ExchangeName]chan spotTypes.SpotWatchEvent
}

func NewSpotWatchlist() spotTypes.SpotWatchlist {
	return &spotWatchlist{
		pairs:        make(map[pairKey]portfolio.Pair),
		baseWatchers: make(map[connector.ExchangeName]chan baseTypes.MarketWatchEvent),
		spotWatchers: make(map[connector.ExchangeName]chan spotTypes.SpotWatchEvent),
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
	w.emitBaseLocked(baseTypes.MarketWatchEvent{
		Requirement: baseTypes.PairRequirement{Exchange: exchange, Pair: pair},
		Type:        baseTypes.PairAdded,
	})
	w.emitSpotLocked(spotTypes.SpotWatchEvent{Exchange: exchange, Pair: pair, Type: spotTypes.SpotPairAdded})
}

func (w *spotWatchlist) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
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
	w.emitSpotLocked(spotTypes.SpotWatchEvent{Exchange: exchange, Pair: pair, Type: spotTypes.SpotPairRemoved})
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

	if ch, ok := w.baseWatchers[exchange]; ok {
		return ch
	}
	ch := make(chan baseTypes.MarketWatchEvent, 128)
	w.baseWatchers[exchange] = ch
	return ch
}

func (w *spotWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.baseWatchers[exchange]; ok {
		delete(w.baseWatchers, exchange)
		close(ch)
	}
}

// SubscribeSpot provides a typed channel for domain-level consumers.
func (w *spotWatchlist) SubscribeSpot(exchange connector.ExchangeName) chan spotTypes.SpotWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.spotWatchers[exchange]; ok {
		return ch
	}
	ch := make(chan spotTypes.SpotWatchEvent, 128)
	w.spotWatchers[exchange] = ch
	return ch
}

func (w *spotWatchlist) UnsubscribeSpot(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.spotWatchers[exchange]; ok {
		delete(w.spotWatchers, exchange)
		close(ch)
	}
}

func (w *spotWatchlist) emitBaseLocked(ev baseTypes.MarketWatchEvent) {
	if ch, ok := w.baseWatchers[ev.Requirement.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

func (w *spotWatchlist) emitSpotLocked(ev spotTypes.SpotWatchEvent) {
	if ch, ok := w.spotWatchers[ev.Exchange]; ok {
		select {
		case ch <- ev:
		default:
		}
	}
}

var _ spotTypes.SpotWatchlist = (*spotWatchlist)(nil)
