package data

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type PairRequirement struct {
	Exchange connector.ExchangeName
	Pair     portfolio.Pair
}

type MarketWatchEventType int

const (
	PairAdded MarketWatchEventType = iota
	PairRemoved
)

type MarketWatchEvent struct {
	Requirement PairRequirement
	Type        MarketWatchEventType
}

// MarketWatchlist manages assets explicitly registered by the application
type MarketWatchlist interface {
	// RequirePair adds a pair to the watchlist, signaling that the ingestor needs data for it.
	RequirePair(exchange connector.ExchangeName, pair portfolio.Pair)
	// ReleasePair removes a pair from the watchlist, signaling that the ingestor can stop providing data for it.
	ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair)

	// GetRequiredPairs Boot time call to get all currently required pairs for an exchange, so the ingestor can start providing data for them.
	GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair

	// Subscribe returns a channel that will receive events whenever pairs are added or removed from the watchlist for the specified exchange.
	Subscribe(exchange connector.ExchangeName) chan MarketWatchEvent

	// Unsubscribe stops sending events to the provided channel for the specified exchange. The channel will be closed after this call.
	Unsubscribe(exchange connector.ExchangeName)
}
