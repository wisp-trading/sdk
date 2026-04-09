package realtime

import (
	"context"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
)

// Ingestor is the realtime ingestor interface for WebSocket connectors.
type Ingestor interface {
	Start(ctx context.Context) error
	Stop() error
}

// Connector provides price feed updates via subscription channels (WebSocket pattern).
type Connector interface {
	// Subscribe returns a channel that receives price updates
	Subscribe(feedID string) <-chan priceFeedTypes.PriceFeedUpdate
	// IsConnected returns whether the connector is active
	IsConnected() bool
	// ErrorChannel returns a channel of errors from the connector
	ErrorChannel() <-chan error
}
