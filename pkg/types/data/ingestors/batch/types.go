package batch

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// BatchIngestorFactory creates batch ingestors dynamically based on registered connectors
type BatchIngestorFactory interface {
	CreateIngestors() []BatchIngestor
}

// BatchIngestor handles REST API data collection for a specific market type
type BatchIngestor interface {
	Start(interval time.Duration) error
	Stop() error
	IsActive() bool
	CollectNow()
	GetMarketType() connector.MarketType
}

// CollectionExtension allows market-specific data collection (funding rates, interest rates, etc.)
type CollectionExtension interface {
	Collect(conn connector.Connector, exchangeName connector.ExchangeName, assets []portfolio.Asset)
}
