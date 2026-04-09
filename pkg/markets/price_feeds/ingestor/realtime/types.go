package realtime

import (
	"context"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
)

// Ingestor consumes price updates from the Pyth connector and persists them to the store.
type Ingestor interface {
	// Start begins listening to the connector and persisting updates
	Start(ctx context.Context) error
	// Stop ceases consuming and closes all subscriber channels
	Stop() error
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
