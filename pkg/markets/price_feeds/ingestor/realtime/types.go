package realtime

import (
	"time"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// Ingestor is the standard batch ingestor interface.
type Ingestor interface {
	Start(interval time.Duration) error
	Stop() error
	IsActive() bool
	CollectNow()
	GetMarketType() connector.MarketType
}

// Connector provides price feed updates via subscription channels.
type Connector interface {
	// Subscribe returns a channel that receives price updates for the given feed
	Subscribe(feedID string) <-chan priceFeedTypes.PriceFeedUpdate
	// IsConnected returns whether the connector is active
	IsConnected() bool
	// ErrorChannel returns a channel of errors from the connector
	ErrorChannel() <-chan error
}
