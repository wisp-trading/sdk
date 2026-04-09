package price_feeds

import (
	"time"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types"
)

// adapter implements types.PriceFeeds by wrapping PriceFeedsStore.
// This decouples the public SDK interface from internal storage implementations.
type adapter struct {
	store priceFeedTypes.PriceFeedsStore
}

// NewPriceFeedsAdapter creates a new adapter from a store.
func NewPriceFeedsAdapter(store priceFeedTypes.PriceFeedsStore) types.PriceFeeds {
	return &adapter{store: store}
}

// GetLatestPrice returns the most recent price for a feed.
func (a *adapter) GetLatestPrice(feedID priceFeedTypes.PriceFeedID) (priceFeedTypes.PriceSnapshot, error) {
	return a.store.GetLatestPrice(feedID)
}

// GetPriceAtTime returns the price closest to a specific time.
func (a *adapter) GetPriceAtTime(feedID priceFeedTypes.PriceFeedID, t time.Time) (priceFeedTypes.PriceSnapshot, error) {
	return a.store.GetPriceAtTime(feedID, t)
}

// GetPriceRange returns all prices within a time range.
func (a *adapter) GetPriceRange(feedID priceFeedTypes.PriceFeedID, start, end time.Time) ([]priceFeedTypes.PriceSnapshot, error) {
	return a.store.GetPriceRange(feedID, start, end)
}

// GetLastUpdated returns the last update time for each feed.
func (a *adapter) GetLastUpdated() map[priceFeedTypes.PriceFeedID]time.Time {
	return a.store.GetLastUpdated()
}
