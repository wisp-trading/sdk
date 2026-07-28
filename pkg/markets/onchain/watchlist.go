package onchain

import (
	baseWatchlist "github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
)

// NewOnchainWatchlist creates an empty onchain-domain watchlist.
func NewOnchainWatchlist() onchainTypes.OnchainWatchlist {
	return baseWatchlist.NewPairWatchlist()
}
