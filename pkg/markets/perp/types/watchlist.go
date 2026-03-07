package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PerpWatchlist is the watchlist for the perp domain.
type PerpWatchlist interface {
	baseTypes.MarketWatchlist
}

// PerpUniverse holds the live set of perp exchanges and their watched pairs.
type PerpUniverse struct {
	Exchanges []connector.Exchange
	Assets    map[connector.ExchangeName][]portfolio.Pair
}

// PerpUniverseProvider computes the live perp trading universe.
type PerpUniverseProvider interface {
	Universe() PerpUniverse
}
