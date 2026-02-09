package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
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
	AccountReader

	// Perp-specific
	GetPerpSymbol(symbol portfolio.Pair) string
}
