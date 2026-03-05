package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// Connector represents a perpetual futures exchange connection.
// AccountReader is the perp-specific override (adds GetMarginBalances) and
// replaces the base connector.AccountReader (GetBalances).
type Connector interface {
	connector.Connector
	connector.MarketDataReader
	connector.OrderExecutor
	FundingRateProvider
	PositionManager
	ContractProvider
	AccountReader

	// Perp-specific
	GetPerpSymbol(symbol portfolio.Pair) string
}
