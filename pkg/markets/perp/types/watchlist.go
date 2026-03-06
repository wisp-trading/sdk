package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
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

// PerpWatchlist embeds the base MarketWatchlist so it satisfies the base ingestor.
// SubscribePerp provides a typed channel for domain-level consumers.
type PerpWatchlist interface {
	baseTypes.MarketWatchlist
	SubscribePerp(exchange connector.ExchangeName) chan PerpWatchEvent
	UnsubscribePerp(exchange connector.ExchangeName)
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
