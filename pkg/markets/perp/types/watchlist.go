package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PerpWatchEventType describes whether a pair was added or removed.
type PerpWatchEventType int

const (
	PerpPairAdded PerpWatchEventType = iota
	PerpPairRemoved
)

// PerpWatchEvent is emitted whenever a pair is added to or removed from the perp watchlist.
type PerpWatchEvent struct {
	Exchange connector.ExchangeName
	Pair     portfolio.Pair
	Type     PerpWatchEventType
}

// PerpWatchlist manages the set of pairs explicitly registered by strategies for perp markets.
// It is the perp-domain equivalent of data.MarketWatchlist.
type PerpWatchlist interface {
	// RequirePair adds a pair to the perp watchlist.
	RequirePair(exchange connector.ExchangeName, pair portfolio.Pair)
	// ReleasePair removes a pair from the perp watchlist.
	ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair)

	// GetRequiredPairs returns all currently required pairs for an exchange.
	GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair

	// Subscribe returns a channel that emits events when pairs are added/removed for an exchange.
	Subscribe(exchange connector.ExchangeName) chan PerpWatchEvent
	// Unsubscribe stops sending events and closes the channel for an exchange.
	Unsubscribe(exchange connector.ExchangeName)
}

// PerpUniverse holds the live set of perp exchanges and their watched pairs.
type PerpUniverse struct {
	// Exchanges are the ready perp connectors.
	Exchanges []connector.Exchange
	// Assets maps each exchange to the pairs currently on the watchlist.
	Assets map[connector.ExchangeName][]portfolio.Pair
}

// PerpUniverseProvider computes the live perp trading universe.
type PerpUniverseProvider interface {
	Universe() PerpUniverse
}
