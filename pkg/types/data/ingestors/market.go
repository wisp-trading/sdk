package ingestors

import (
	"context"
)

// MarketDataCoordinator manages data ingestion across all market types
// Orchestrates multiple RealtimeIngestors and BatchIngestors
type MarketDataCoordinator interface {
	IsRunning() bool
	StartDataCollection(ctx context.Context) error
	StopDataCollection() error
	GetStatus() map[string]interface{}
	ForceCollectNow()
	RestartRealtime(ctx context.Context) error
}
