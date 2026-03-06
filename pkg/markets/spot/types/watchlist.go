package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
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

// SpotWatchlist embeds the base MarketWatchlist so it satisfies the base ingestor.
// SubscribeSpot provides a typed channel for domain-level consumers.
type SpotWatchlist interface {
	baseTypes.MarketWatchlist
	SubscribeSpot(exchange connector.ExchangeName) chan SpotWatchEvent
	UnsubscribeSpot(exchange connector.ExchangeName)
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
