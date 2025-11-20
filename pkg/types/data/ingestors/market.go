package ingestors

import (
	"context"
	"time"
)

// MarketDataCoordinator manages hybrid data ingestion (realtime + batch)
type MarketDataCoordinator interface {
	IsRunning() bool
	StartDataCollection(ctx context.Context) error
	StopDataCollection() error
	GetStatus() map[string]interface{}
	ForceCollectNow()
	RestartRealtime(ctx context.Context) error
}

// RealtimeIngestor handles real-time market data ingestion via WebSocket
type RealtimeIngestor interface {
	Start(ctx context.Context) error
	Stop() error
	IsActive() bool
}

// BatchIngestor handles periodic batch market data collection via REST
type BatchIngestor interface {
	Start(interval time.Duration) error
	Stop() error
	IsActive() bool
	CollectNow()
}
