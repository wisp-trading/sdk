package prediction

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// internal key for the map.
type predictionKey struct {
	Exchange connector.ExchangeName
	MarketID predictionconnector.MarketID
}

// predictionWatchlist is the concrete implementation of PredictionWatchlist.
type predictionWatchlist struct {
	mu sync.RWMutex

	markets map[predictionKey]predictionconnector.Market

	// exactly one channel per exchange
	watchers map[connector.ExchangeName]chan types.PredictionWatchEvent
}

func NewPredictionWatchlist() types.PredictionWatchlist {
	return &predictionWatchlist{
		markets:  make(map[predictionKey]predictionconnector.Market),
		watchers: make(map[connector.ExchangeName]chan types.PredictionWatchEvent),
	}
}

func (w *predictionWatchlist) RequireMarket(exchange connector.ExchangeName, market predictionconnector.Market) {
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

	w.emitEventLocked(types.PredictionWatchEvent{
		Exchange: exchange,
		Market:   market,
		Type:     types.PredictionMarketAdded,
	})
}

func (w *predictionWatchlist) ReleaseMarket(exchange connector.ExchangeName, marketID predictionconnector.MarketID) {
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

	w.emitEventLocked(types.PredictionWatchEvent{
		Exchange: exchange,
		Market:   market,
		Type:     types.PredictionMarketRemoved,
	})
}

func (w *predictionWatchlist) GetRequiredMarkets(exchange connector.ExchangeName) []predictionconnector.Market {
	w.mu.RLock()
	defer w.mu.RUnlock()

	out := make([]predictionconnector.Market, 0, len(w.markets))
	for key, m := range w.markets {
		if key.Exchange == exchange {
			out = append(out, m)
		}
	}

	return out
}

func (w *predictionWatchlist) GetAllMarkets() map[connector.ExchangeName][]predictionconnector.Market {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make(map[connector.ExchangeName][]predictionconnector.Market)
	for key, m := range w.markets {
		result[key.Exchange] = append(result[key.Exchange], m)
	}
	return result
}

func (w *predictionWatchlist) Subscribe(exchange connector.ExchangeName) chan types.PredictionWatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		// Already have a channel for this exchange; reuse it.
		return ch
	}

	ch := make(chan types.PredictionWatchEvent, 128)
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

func (w *predictionWatchlist) emitEventLocked(ev types.PredictionWatchEvent) {
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
