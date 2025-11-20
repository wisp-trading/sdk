package ingestors

import (
	"context"
	"time"
)

// MarketDataIngestor is the base interface for all market data ingestors
type MarketDataIngestor interface {
	// Start begins the ingestion process
	Start(ctx context.Context) error

	// Stop halts the ingestion process
	Stop() error

	// IsActive returns whether the ingestor is currently running
	IsActive() bool
}

// RealtimeIngestor handles real-time market data ingestion via WebSocket
type RealtimeIngestor interface {
	MarketDataIngestor
}

// BatchIngestor handles periodic batch market data collection via REST
type BatchIngestor interface {
	MarketDataIngestor

	// StartWithInterval begins batch collection with the specified interval
	StartWithInterval(interval time.Duration) error

	// CollectNow triggers an immediate collection
	CollectNow()
}

// MarketDataCoordinator manages hybrid data ingestion (realtime + batch)
type MarketDataCoordinator interface {
	// IsRunning returns whether the coordinator is active
	IsRunning() bool

	// StartDataCollection starts both realtime and batch ingestion
	StartDataCollection(ctx context.Context) error

	// StopDataCollection stops all ingestion
	StopDataCollection() error

	// GetStatus returns the current status of all ingestors
	GetStatus() map[string]interface{}

	// ForceCollectNow triggers immediate batch collection
	ForceCollectNow()

	// RestartRealtime restarts the realtime ingestor
	RestartRealtime(ctx context.Context) error
}

// PositionCoordinator handles trade backfill on startup
type PositionCoordinator interface {
	// Start begins the coordinator and performs initial backfill
	Start(ctx context.Context) error

	// Stop halts the coordinator
	Stop() error

	// IsActive returns whether the coordinator is running
	IsActive() bool

	// GetStatus returns the current status
	GetStatus() map[string]interface{}
}
