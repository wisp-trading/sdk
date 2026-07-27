package perp

import (
	baseWatchlist "github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
)

// NewPerpWatchlist creates an empty perp-domain watchlist (shared pair impl).
func NewPerpWatchlist() perpTypes.PerpWatchlist {
	return baseWatchlist.NewPairWatchlist()
}
