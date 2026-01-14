package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/common"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Connector represents a perpetual futures exchange connection
type Connector interface {
	common.BaseConnector
	common.MarketDataReader
	common.OrderExecutor
	common.AccountReader
	FundingRateProvider
	PositionManager
	ContractProvider

	// Perp-specific
	FetchAvailableAssets() ([]portfolio.Asset, error)
	GetPerpSymbol(symbol portfolio.Asset) string
}
