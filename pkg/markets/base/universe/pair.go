// Package universe builds live exchange/asset sets for pair-based domains (spot/perp).
package universe

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PairUniverse is exchanges + watched pairs (shared shape for spot/perp).
type PairUniverse struct {
	Exchanges []connector.Exchange
	Assets    map[connector.ExchangeName][]portfolio.Pair
}

// BuildPairUniverse constructs a live universe from ready connectors and a watchlist.
// readyNames are exchange names already filtered ready for the domain.
func BuildPairUniverse(
	readyNames []connector.ExchangeName,
	marketType connector.MarketType,
	watchlist baseTypes.MarketWatchlist,
) PairUniverse {
	exchanges := make([]connector.Exchange, 0, len(readyNames))
	assets := make(map[connector.ExchangeName][]portfolio.Pair)

	for _, name := range readyNames {
		exchanges = append(exchanges, connector.Exchange{
			Name:       name,
			MarketType: marketType,
		})
		if pairs := watchlist.GetRequiredPairs(name); len(pairs) > 0 {
			assets[name] = pairs
		}
	}

	return PairUniverse{Exchanges: exchanges, Assets: assets}
}
