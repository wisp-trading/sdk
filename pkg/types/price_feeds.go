package types

import (
	"time"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
)

// PriceFeeds is the public SDK interface for accessing price feed data.
// Strategies use this to query price feeds without coupling to internal storage implementations.
type PriceFeeds interface {
	// GetLatestPrice returns the most recent price for a feed.
	GetLatestPrice(feedID priceFeedTypes.PriceFeedID) (priceFeedTypes.PriceSnapshot, error)

	// GetPriceAtTime returns the price closest to a specific time.
	GetPriceAtTime(feedID priceFeedTypes.PriceFeedID, t time.Time) (priceFeedTypes.PriceSnapshot, error)

	// GetPriceRange returns all prices within a time range.
	GetPriceRange(feedID priceFeedTypes.PriceFeedID, start, end time.Time) ([]priceFeedTypes.PriceSnapshot, error)

	// GetLastUpdated returns the last update time for each feed.
	GetLastUpdated() map[priceFeedTypes.PriceFeedID]time.Time
}
