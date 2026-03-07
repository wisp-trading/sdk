package config

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PerpDomainConfig is the perp domain's decoded view of the strategy config.
// It is produced once at startup from the raw StartupConfig, filtered to perp
// exchanges only, and stored in the fx graph so every perp component can depend
// on it directly without cross-domain coupling.
type PerpDomainConfig struct {
	// Assets maps each perp exchange to the pairs it should watch.
	Assets map[connector.ExchangeName][]portfolio.Pair
}
