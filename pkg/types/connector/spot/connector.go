package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Connector represents a spot market exchange connection
type Connector interface {
	connector.Connector
	connector.MarketDataReader
	connector.OrderExecutor
	connector.AccountReader

	// Spot-specific
	FetchAvailableAssets() ([]portfolio.Asset, error)
}
