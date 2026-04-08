package types

import (
	"time"

	optionsConnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// OptionsUniverse represents the available options trading universe
type OptionsUniverse struct {
	Exchanges []connector.Exchange
	// Map of exchange -> pair -> expirations
	Expirations map[connector.ExchangeName]map[portfolio.Pair][]time.Time
	// Map of exchange -> pair -> expiration -> strikes
	Strikes map[connector.ExchangeName]map[portfolio.Pair]map[time.Time][]float64
}

// Options-specific type aliases for clarity
type (
	OptionContract = optionsConnector.OptionContract
	OptionData     = optionsConnector.OptionData
	Greeks         = optionsConnector.Greeks
)
