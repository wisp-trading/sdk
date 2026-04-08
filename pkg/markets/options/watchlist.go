package options

import (
	"sync"
	"time"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// expirationKey uniquely identifies a slow-watched expiration
type expirationKey struct {
	Exchange   connector.ExchangeName
	Pair       portfolio.Pair
	Expiration time.Time
}

// instrumentKey uniquely identifies a fast-watched contract
type instrumentKey struct {
	Exchange   connector.ExchangeName
	Pair       portfolio.Pair
	Expiration time.Time
	Strike     float64
	OptionType string
}

// optionsWatchlist is the concrete implementation of OptionsWatchlist
type optionsWatchlist struct {
	mu sync.RWMutex

	// Slow watch — full expiration polling via REST
	expirations map[expirationKey]bool
	strikes     map[expirationKey][]float64

	// Fast watch — real-time order book per specific contract
	instruments map[instrumentKey]optionsTypes.OptionContract

	// Event subscribers per exchange
	watchers map[connector.ExchangeName]chan optionsTypes.WatchEvent
}

// NewOptionsWatchlist creates a new options watchlist
func NewOptionsWatchlist() optionsTypes.OptionsWatchlist {
	return &optionsWatchlist{
		expirations: make(map[expirationKey]bool),
		strikes:     make(map[expirationKey][]float64),
		instruments: make(map[instrumentKey]optionsTypes.OptionContract),
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
		Type:       optionsTypes.WatchEventExpirationAdded,
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
		Type:       optionsTypes.WatchEventExpirationRemoved,
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

// WatchInstrument adds a specific contract to the fast watch tier.
// The realtime ingestor will open a WebSocket order book subscription for it.
func (w *optionsWatchlist) WatchInstrument(
	exchange connector.ExchangeName,
	contract optionsTypes.OptionContract,
) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := instrumentKey{
		Exchange:   exchange,
		Pair:       contract.Pair,
		Expiration: contract.Expiration,
		Strike:     contract.Strike,
		OptionType: contract.OptionType,
	}

	if _, exists := w.instruments[key]; exists {
		return
	}

	w.instruments[key] = contract

	w.emitLocked(optionsTypes.WatchEvent{
		Type:       optionsTypes.WatchEventInstrumentWatched,
		Exchange:   exchange,
		Pair:       contract.Pair,
		Expiration: contract.Expiration,
		Contract:   &contract,
	})
}

// UnwatchInstrument removes a contract from the fast watch tier.
func (w *optionsWatchlist) UnwatchInstrument(
	exchange connector.ExchangeName,
	contract optionsTypes.OptionContract,
) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := instrumentKey{
		Exchange:   exchange,
		Pair:       contract.Pair,
		Expiration: contract.Expiration,
		Strike:     contract.Strike,
		OptionType: contract.OptionType,
	}

	if _, exists := w.instruments[key]; !exists {
		return
	}

	delete(w.instruments, key)

	w.emitLocked(optionsTypes.WatchEvent{
		Type:       optionsTypes.WatchEventInstrumentUnwatched,
		Exchange:   exchange,
		Pair:       contract.Pair,
		Expiration: contract.Expiration,
		Contract:   &contract,
	})
}

// GetWatchedInstruments returns all fast-watched contracts for an exchange.
func (w *optionsWatchlist) GetWatchedInstruments(
	exchange connector.ExchangeName,
) []optionsTypes.OptionContract {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var contracts []optionsTypes.OptionContract
	for key, contract := range w.instruments {
		if key.Exchange == exchange {
			contracts = append(contracts, contract)
		}
	}
	return contracts
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
		Type:       optionsTypes.WatchEventStrikesUpdated,
		Exchange:   exchange,
		Pair:       pair,
		Expiration: expiration,
		Strikes:    strikes,
	})
}

// Ensure optionsWatchlist implements OptionsWatchlist
var _ optionsTypes.OptionsWatchlist = (*optionsWatchlist)(nil)
