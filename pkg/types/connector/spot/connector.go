package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/common"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Connector represents a spot market exchange connection
type Connector interface {
	common.BaseConnector
	common.MarketDataReader
	common.OrderExecutor
	common.AccountReader

	// Spot-specific
	FetchAvailableAssets() ([]portfolio.Asset, error)
}
