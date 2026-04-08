package options

import (
	"sync"
	"time"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// expirationKey uniquely identifies an expiration to watch
type expirationKey struct {
	Exchange   connector.ExchangeName
	Pair       portfolio.Pair
	Expiration time.Time
}

// optionsWatchlist is the concrete implementation of OptionsWatchlist
type optionsWatchlist struct {
	mu sync.RWMutex

	// expirations: map[expirationKey]bool - tracks what we're watching
	expirations map[expirationKey]bool

	// strikes: map[expirationKey][]float64 - discovered strikes per expiration
	strikes map[expirationKey][]float64

	// watchers: map[ExchangeName]<-chan WatchEvent
	watchers map[connector.ExchangeName]chan optionsTypes.WatchEvent

	// connector registry for on-demand strike discovery
	connectorRegistry interface {
		FilterOptions(opts interface{}) []interface{}
	}
}

// NewOptionsWatchlist creates a new options watchlist
func NewOptionsWatchlist() optionsTypes.OptionsWatchlist {
	return &optionsWatchlist{
		expirations: make(map[expirationKey]bool),
		strikes:     make(map[expirationKey][]float64),
		watchers:    make(map[connector.ExchangeName]chan optionsTypes.WatchEvent),
	}
}

// RequireExpiration adds an expiration to the watchlist
func (w *optionsWatchlist) RequireExpiration(
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	expiration time.Time,
) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := expirationKey{Exchange: exchange, Pair: pair, Expiration: expiration}

	// Already watching this expiration
	if w.expirations[key] {
		return nil
	}

	// Mark as watched
	w.expirations[key] = true

	// Emit event
	w.emitLocked(optionsTypes.WatchEvent{
		Type:       "ExpirationAdded",
		Exchange:   exchange,
		Pair:       pair,
		Expiration: expiration,
	})

	return nil
}

// ReleaseExpiration removes an expiration from the watchlist
func (w *optionsWatchlist) ReleaseExpiration(
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	expiration time.Time,
) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := expirationKey{Exchange: exchange, Pair: pair, Expiration: expiration}

	if !w.expirations[key] {
		return nil
	}

	delete(w.expirations, key)
	delete(w.strikes, key)

	w.emitLocked(optionsTypes.WatchEvent{
		Type:       "ExpirationRemoved",
		Exchange:   exchange,
		Pair:       pair,
		Expiration: expiration,
	})

	return nil
}

// GetAvailableStrikes returns discovered strikes for a watched expiration
func (w *optionsWatchlist) GetAvailableStrikes(
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	expiration time.Time,
) []float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	key := expirationKey{Exchange: exchange, Pair: pair, Expiration: expiration}
	strikes, ok := w.strikes[key]
	if !ok {
		return []float64{}
	}

	return strikes
}

// GetWatchedExpirations returns all watched expirations for an exchange
func (w *optionsWatchlist) GetWatchedExpirations(
	exchange connector.ExchangeName,
) map[portfolio.Pair][]time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make(map[portfolio.Pair][]time.Time)

	for key := range w.expirations {
		if key.Exchange != exchange {
			continue
		}
		result[key.Pair] = append(result[key.Pair], key.Expiration)
	}

	return result
}

// Subscribe returns a channel for watchlist updates
func (w *optionsWatchlist) Subscribe(exchange connector.ExchangeName) <-chan optionsTypes.WatchEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		return ch
	}

	ch := make(chan optionsTypes.WatchEvent, 128)
	w.watchers[exchange] = ch
	return ch
}

// Unsubscribe closes the watchlist channel for an exchange
func (w *optionsWatchlist) Unsubscribe(exchange connector.ExchangeName) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, ok := w.watchers[exchange]; ok {
		delete(w.watchers, exchange)
		close(ch)
	}
}

// emitLocked sends an event to the watcher channel (must hold lock)
func (w *optionsWatchlist) emitLocked(ev optionsTypes.WatchEvent) {
	if ch, ok := w.watchers[ev.Exchange]; ok {
		select {
		case ch <- ev:
		default:
			// Drop if channel is full
		}
	}
}

// SetStrikes updates the discovered strikes for an expiration (internal use)
func (w *optionsWatchlist) SetStrikes(
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	expiration time.Time,
	strikes []float64,
) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := expirationKey{Exchange: exchange, Pair: pair, Expiration: expiration}
	w.strikes[key] = strikes

	w.emitLocked(optionsTypes.WatchEvent{
		Type:       "StrikesUpdated",
		Exchange:   exchange,
		Pair:       pair,
		Expiration: expiration,
		Strikes:    strikes,
	})
}

// Ensure optionsWatchlist implements OptionsWatchlist
var _ optionsTypes.OptionsWatchlist = (*optionsWatchlist)(nil)
