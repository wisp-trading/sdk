package spot

import (
	baseWatchlist "github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
)

// NewSpotWatchlist creates an empty spot-domain watchlist (shared pair impl).
func NewSpotWatchlist() spotTypes.SpotWatchlist {
	return baseWatchlist.NewPairWatchlist()
}
