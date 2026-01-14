package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Connector represents a perpetual futures exchange connection
type Connector interface {
	connector.Connector
	connector.MarketDataReader
	connector.OrderExecutor
	connector.AccountReader
	FundingRateProvider
	PositionManager
	ContractProvider

	// Perp-specific
	FetchAvailableAssets() ([]portfolio.Asset, error)
	GetPerpSymbol(symbol portfolio.Asset) string
}
