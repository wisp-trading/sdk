package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// SpotWatchlist is the watchlist for the spot domain.
type SpotWatchlist interface {
	baseTypes.MarketWatchlist
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
