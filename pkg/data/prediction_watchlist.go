package data

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	prediction "github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/data"
)

// internal key for the map.
type predictionKey struct {
	Exchange connector.ExchangeName
	MarketID prediction.MarketID
}

// predictionWatchlist is the concrete implementation of PredictionWatchlist.
type predictionWatchlist struct {
	mu sync.RWMutex

	markets map[predictionKey]prediction.Market

	// exactly one channel per exchange
	watchers map[connector.ExchangeName]chan data.PredictionWatchEvent
}

func NewPredictionWatchlist() data.PredictionWatchlist {
	return &predictionWatchlist{
		markets:  make(map[predictionKey]prediction.Market),
		watchers: make(map[connector.ExchangeName]chan data.PredictionWatchEvent),
	}
}

func (w *predictionWatchlist) RequireMarket(exchange connector.ExchangeName, market prediction.Market) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := predictionKey{
		Exchange: exchange,
		MarketID: market.MarketID,
	}

	if _, exists := w.markets[key]; exists {
		// Already required; no-op.
		return
	}

	w.markets[key] = market

	w.emitEventLocked(data.PredictionWatchEvent{
		Exchange: exchange,
		Market:   market,
		Type:     data.PredictionMarketAdded,
	})
}

func (w *predictionWatchlist) ReleaseMarket(exchange connector.ExchangeName, marketID prediction.MarketID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := predictionKey{
		Exchange: exchange,
		MarketID: marketID,
	}

	market, exists := w.markets[key]
	if !exists {
		// Not present; no-op.
		return
	}

	delete(w.markets, key)

	w.emitEventLocked(data.PredictionWatchEvent{
		Exchange: exchange,
		Market:   market,
		Type:     data.PredictionMarketRemoved,
	})
}

func (w *predictionWatchlist) GetRequiredMarkets(exchange connector.ExchangeName) []prediction.Market {
	w.mu.RLock()
	defer w.mu.RUnlock()

	out := make([]prediction.Market, 0, len(w.markets))
	for key, m := range w.markets {
		if key.Exchange == exchange {
			out = append(out, m)
		}
	}

	return out
}

func (w *predictionWatchlist) Subscribe(exchange connector.ExchangeName) chan data.PredictionWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		// Already have a channel for this exchange; reuse it.
		return ch
	}

	ch := make(chan data.PredictionWatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

func (w *predictionWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch, ok := w.watchers[exchange]
	if !ok {
		return
	}

	delete(w.watchers, exchange)
	close(ch)
}

func (w *predictionWatchlist) emitEventLocked(ev data.PredictionWatchEvent) {
	ch, ok := w.watchers[ev.Exchange]
	if !ok {
		return
	}

	select {
	case ch <- ev:
	default:
		// drop or log slow watcher
	}
}
