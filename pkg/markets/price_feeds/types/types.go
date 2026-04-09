package types

import (
	"time"

	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// PriceFeedID identifies a price feed (e.g., "pyth:wtik6", "chainlink:btc")
type PriceFeedID string

// PriceSnapshot represents a single price update
type PriceSnapshot struct {
	FeedID    PriceFeedID
	Price     float64
	Timestamp time.Time
}

// PriceFeedsStore is the store for price feed history
type PriceFeedsStore interface {
	// Price history operations
	RecordPrice(snap PriceSnapshot) error
	GetLatestPrice(feedID PriceFeedID) (PriceSnapshot, error)
	GetPriceAtTime(feedID PriceFeedID, t time.Time) (PriceSnapshot, error)
	GetPriceRange(feedID PriceFeedID, start, end time.Time) ([]PriceSnapshot, error)

	// Cleanup
	PruneOldData(olderThan time.Time) error

	// Metadata
	GetLastUpdated() map[PriceFeedID]time.Time
	UpdateLastUpdated(feedID PriceFeedID)
}

// PriceFeedUpdate is what ingestors emit
type PriceFeedUpdate struct {
	FeedID    PriceFeedID
	Price     float64
	Timestamp time.Time
	Source    connector.ExchangeName // e.g., "pyth", "chainlink"
}

// PriceFeedsBatchIngestorFactory creates batch ingestors for price feeds.
type PriceFeedsBatchIngestorFactory interface {
	CreateIngestors() []batchTypes.BatchIngestor
}
