package store

import (
	"fmt"
	"sync"
	"time"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
)

type priceFeedsStore struct {
	mu sync.RWMutex

	// Price history: feedID -> time-ordered snapshots
	prices map[priceFeedTypes.PriceFeedID][]priceFeedTypes.PriceSnapshot

	// Metadata: feedID -> last update time
	lastUpdated map[priceFeedTypes.PriceFeedID]time.Time
}

// NewStore creates a new price feeds store
func NewStore() priceFeedTypes.PriceFeedsStore {
	return &priceFeedsStore{
		prices:      make(map[priceFeedTypes.PriceFeedID][]priceFeedTypes.PriceSnapshot),
		lastUpdated: make(map[priceFeedTypes.PriceFeedID]time.Time),
	}
}

// RecordPrice stores a price update
func (s *priceFeedsStore) RecordPrice(snap priceFeedTypes.PriceSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.prices[snap.FeedID]; !ok {
		s.prices[snap.FeedID] = []priceFeedTypes.PriceSnapshot{}
	}

	s.prices[snap.FeedID] = append(s.prices[snap.FeedID], snap)

	// Update metadata
	s.lastUpdated[snap.FeedID] = snap.Timestamp

	return nil
}

// GetLatestPrice returns the most recent price for a feed
func (s *priceFeedsStore) GetLatestPrice(feedID priceFeedTypes.PriceFeedID) (priceFeedTypes.PriceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots, ok := s.prices[feedID]
	if !ok || len(snapshots) == 0 {
		return priceFeedTypes.PriceSnapshot{}, fmt.Errorf("no price for feed %s", feedID)
	}

	return snapshots[len(snapshots)-1], nil
}

// GetPriceAtTime returns the price at or before time t
func (s *priceFeedsStore) GetPriceAtTime(feedID priceFeedTypes.PriceFeedID, t time.Time) (priceFeedTypes.PriceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots, ok := s.prices[feedID]
	if !ok || len(snapshots) == 0 {
		return priceFeedTypes.PriceSnapshot{}, fmt.Errorf("no price for feed %s", feedID)
	}

	// Binary search for price at or before t
	for i := len(snapshots) - 1; i >= 0; i-- {
		if snapshots[i].Timestamp.Before(t) || snapshots[i].Timestamp.Equal(t) {
			return snapshots[i], nil
		}
	}

	return priceFeedTypes.PriceSnapshot{}, fmt.Errorf("no price for feed %s before %v", feedID, t)
}

// GetPriceRange returns all prices between start and end times
func (s *priceFeedsStore) GetPriceRange(feedID priceFeedTypes.PriceFeedID, start, end time.Time) ([]priceFeedTypes.PriceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots, ok := s.prices[feedID]
	if !ok {
		return nil, fmt.Errorf("no prices for feed %s", feedID)
	}

	var result []priceFeedTypes.PriceSnapshot
	for _, snap := range snapshots {
		if (snap.Timestamp.Equal(start) || snap.Timestamp.After(start)) &&
			(snap.Timestamp.Before(end) || snap.Timestamp.Equal(end)) {
			result = append(result, snap)
		}
	}

	return result, nil
}

// PruneOldData removes prices older than the specified time
func (s *priceFeedsStore) PruneOldData(olderThan time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for feedID, snapshots := range s.prices {
		pruned := 0
		for i, snap := range snapshots {
			if snap.Timestamp.After(olderThan) {
				pruned = i
				break
			}
		}
		if pruned > 0 {
			s.prices[feedID] = snapshots[pruned:]
		}
	}

	return nil
}

// GetLastUpdated returns the last updated map
func (s *priceFeedsStore) GetLastUpdated() map[priceFeedTypes.PriceFeedID]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[priceFeedTypes.PriceFeedID]time.Time)
	for k, v := range s.lastUpdated {
		result[k] = v
	}
	return result
}

// UpdateLastUpdated updates the last updated time
func (s *priceFeedsStore) UpdateLastUpdated(feedID priceFeedTypes.PriceFeedID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastUpdated[feedID] = time.Now()
}

// Ensure priceFeedsStore implements PriceFeedsStore
var _ priceFeedTypes.PriceFeedsStore = (*priceFeedsStore)(nil)
