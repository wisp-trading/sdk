package types

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// OptionsWatchlist manages two tiers of options market observation:
//
// Slow watch — full expiration polling via REST every ~30s.
// Tracks all strikes for a given expiration. Used for scanning
// the option chain, comparing funding income vs theta cost, etc.
//
// Fast watch — real-time WebSocket order book for a specific contract.
// Used when assessing liquidity before entry, or actively managing a position.
type OptionsWatchlist interface {

	// --- Slow watch (batch / REST) ---

	// RequireExpiration adds an expiration to the slow watch.
	// The batch ingestor will poll all strikes for it on every tick.
	RequireExpiration(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) error
	ReleaseExpiration(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) error

	// SetStrikes records discovered strikes for an expiration (called by batch ingestor).
	SetStrikes(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time, strikes []float64)

	// GetAvailableStrikes returns the known strikes for a slow-watched expiration.
	GetAvailableStrikes(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) []float64

	// GetWatchedExpirations returns all slow-watched expirations for an exchange.
	GetWatchedExpirations(exchange connector.ExchangeName) map[portfolio.Pair][]time.Time

	// --- Fast watch (WebSocket / real-time order book) ---

	// WatchInstrument adds a specific contract to the fast watch.
	// Triggers a WebSocket order book subscription for that instrument.
	WatchInstrument(exchange connector.ExchangeName, contract OptionContract)

	// UnwatchInstrument removes a contract from the fast watch.
	UnwatchInstrument(exchange connector.ExchangeName, contract OptionContract)

	// GetWatchedInstruments returns all fast-watched contracts for an exchange.
	GetWatchedInstruments(exchange connector.ExchangeName) []OptionContract

	// Subscribe returns a channel of watchlist events for an exchange.
	Subscribe(exchange connector.ExchangeName) <-chan WatchEvent
	// Unsubscribe closes the event channel for an exchange.
	Unsubscribe(exchange connector.ExchangeName)
}

// Watch event type constants — use these instead of bare strings.
const (
	WatchEventExpirationAdded    = "ExpirationAdded"
	WatchEventExpirationRemoved  = "ExpirationRemoved"
	WatchEventStrikesUpdated     = "StrikesUpdated"
	WatchEventInstrumentWatched  = "InstrumentWatched"
	WatchEventInstrumentUnwatched = "InstrumentUnwatched"
)

// WatchEvent represents a change in the options watchlist.
type WatchEvent struct {
	Type       string // One of the WatchEvent* constants above.
	Exchange   connector.ExchangeName
	Pair       portfolio.Pair
	Expiration time.Time
	Strikes    []float64       // Populated for WatchEventStrikesUpdated
	Contract   *OptionContract // Populated for WatchEventInstrumentWatched / WatchEventInstrumentUnwatched
}
