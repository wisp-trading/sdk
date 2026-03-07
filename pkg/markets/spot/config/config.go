package config

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// SpotDomainConfig is the spot domain's decoded view of the strategy config.
// It is produced once at startup from the raw StartupConfig, filtered to spot
// exchanges only, and stored in the fx graph so every spot component can depend
// on it directly without cross-domain coupling.
type SpotDomainConfig struct {
	// Assets maps each spot exchange to the pairs it should watch.
	Assets map[connector.ExchangeName][]portfolio.Pair
}
