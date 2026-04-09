package batch

import (
	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
)

// Connector provides price feed data via REST-like polling methods.
type Connector interface {
	// GetLatestPrices fetches the latest prices for all configured feeds
	GetLatestPrices() ([]priceFeedTypes.PriceFeedUpdate, error)
	// GetConnectorInfo returns information about this connector
	GetConnectorInfo() ConnectorInfo
}

type ConnectorInfo struct {
	Name string
}
