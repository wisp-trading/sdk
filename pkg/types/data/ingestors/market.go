package ingestors

import (
	"context"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
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

// RealtimeIngestor handles WebSocket data ingestion for a specific market type
type RealtimeIngestor interface {
	Start(ctx context.Context) error
	Stop() error
	IsActive() bool
	GetMarketType() connector.MarketType
	GetActiveConnections() map[connector.ExchangeName]interface{} // Returns type-specific connectors
}

// BatchIngestor handles REST API data collection for a specific market type
type BatchIngestor interface {
	Start(interval time.Duration) error
	Stop() error
	IsActive() bool
	CollectNow()
	GetMarketType() connector.MarketType
}
