package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type SpotWatchEventType int

const (
	SpotPairAdded SpotWatchEventType = iota
	SpotPairRemoved
)

type SpotWatchEvent struct {
	Exchange connector.ExchangeName
	Pair     portfolio.Pair
	Type     SpotWatchEventType
}

// SpotWatchlist manages the set of pairs registered by strategies for spot markets.
type SpotWatchlist interface {
	RequirePair(exchange connector.ExchangeName, pair portfolio.Pair)
	ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair)
	GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair
	Subscribe(exchange connector.ExchangeName) chan SpotWatchEvent
	Unsubscribe(exchange connector.ExchangeName)
}

// SpotUniverse holds the live set of spot exchanges and their watched pairs.
type SpotUniverse struct {
	Exchanges []connector.Exchange
	Assets    map[connector.ExchangeName][]portfolio.Pair
}

// SpotUniverseProvider computes the live spot trading universe.
type SpotUniverseProvider interface {
	Universe() SpotUniverse
}
