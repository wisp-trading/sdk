package prediction

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type universeProvider struct {
	watchlist         types.PredictionWatchlist
	connectorRegistry registry.ConnectorRegistry
}

func NewPredictionUniverseProvider(
	watchlist types.PredictionWatchlist,
	connectorRegistry registry.ConnectorRegistry,
) types.PredictionUniverseProvider {
	return &universeProvider{
		watchlist:         watchlist,
		connectorRegistry: connectorRegistry,
	}
}

// Universe returns the live prediction trading universe — always current, never cached.
func (u *universeProvider) Universe() types.PredictionUniverse {
	readyConnectors := u.connectorRegistry.FilterPrediction(
		registry.NewFilter().ReadyOnly().Build(),
	)

	exchanges := make([]connector.Exchange, 0, len(readyConnectors))
	for _, conn := range readyConnectors {
		info := conn.GetConnectorInfo()
		exchanges = append(exchanges, connector.Exchange{
			Name:       info.Name,
			MarketType: connector.MarketTypePrediction,
		})
	}

	return types.PredictionUniverse{
		Exchanges: exchanges,
		Markets:   u.watchlist.GetAllMarkets(),
	}
}
