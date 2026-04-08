package types

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// OptionsWatchlist manages the set of options expirations the SDK should track
type OptionsWatchlist interface {

	// Require an expiration (auto-discovers strikes on-demand)
	RequireExpiration(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) error
	ReleaseExpiration(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) error

	// Get strikes for watched expiration
	GetAvailableStrikes(exchange connector.ExchangeName, pair portfolio.Pair, expiration time.Time) []float64

	// Get all watched expirations per exchange
	// Returns map[pair][]time.Time
	GetWatchedExpirations(exchange connector.ExchangeName) map[portfolio.Pair][]time.Time
}

// WatchEvent represents a change in the options watchlist
type WatchEvent struct {
	Type       string // "ExpirationAdded", "ExpirationRemoved", "StrikesUpdated"
	Exchange   connector.ExchangeName
	Pair       portfolio.Pair
	Expiration time.Time
	Strikes    []float64 // Only for StrikesUpdated
}
