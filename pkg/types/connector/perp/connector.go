package perp

import (
	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
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
	GetPerpSymbol(symbol portfolio.Asset) string
}
