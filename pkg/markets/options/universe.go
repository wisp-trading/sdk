package options

import (
	"time"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type universeProvider struct {
	watchlist         optionsTypes.OptionsWatchlist
	connectorRegistry registry.ConnectorRegistry
}

// NewOptionsUniverseProvider creates a universe provider for the options domain
func NewOptionsUniverseProvider(
	watchlist optionsTypes.OptionsWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) optionsTypes.OptionsUniverseProvider {
	return &universeProvider{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the live options trading universe — always current, never cached
func (u *universeProvider) Universe() optionsTypes.OptionsUniverse {
	readyConnectors := u.connectorRegistry.FilterOptions(
		registry.NewFilter().ReadyOnly().Build(),
	)

	exchanges := make([]connector.Exchange, 0, len(readyConnectors))
	expirations := make(map[connector.ExchangeName]map[portfolio.Pair][]time.Time)
	strikes := make(map[connector.ExchangeName]map[portfolio.Pair]map[time.Time][]float64)

	for _, conn := range readyConnectors {
		info := conn.GetConnectorInfo()

		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypeOptions,
		})

		// Get watched expirations for this exchange
		watchedExpirations := u.watchlist.GetWatchedExpirations(info.Name)
		if len(watchedExpirations) == 0 {
			continue
		}

		// Add expirations and strikes for this exchange
		expirations[info.Name] = watchedExpirations
		strikes[info.Name] = make(map[portfolio.Pair]map[time.Time][]float64)

		for pair, expTimes := range watchedExpirations {
			strikes[info.Name][pair] = make(map[time.Time][]float64)
			for _, expTime := range expTimes {
				strikeList := u.watchlist.GetAvailableStrikes(info.Name, pair, expTime)
				strikes[info.Name][pair][expTime] = strikeList
			}
		}
	}

	return optionsTypes.OptionsUniverse{
		Exchanges:   exchanges,
		Expirations: expirations,
		Strikes:     strikes,
	}
}

var _ optionsTypes.OptionsUniverseProvider = (*universeProvider)(nil)
