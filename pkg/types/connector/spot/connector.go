package spot

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// Connector represents a spot market exchange connection
type Connector interface {
	connector.Connector
	connector.MarketDataReader
	connector.OrderExecutor
	connector.AccountReader
}
