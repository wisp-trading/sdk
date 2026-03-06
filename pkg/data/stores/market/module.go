package market

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/spot"
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"go.uber.org/fx"
)

var Module = fx.Options(
	spot.Module,
	fx.Provide(
		NewMarketRegistry,
	),
	fx.Invoke(
		registerStores,
	),
)

// registerStores registers the spot store with the registry.
// Perp store is registered by pkg/markets/perp/module.go.
// Prediction store is registered by pkg/markets/prediction/module.go.
func registerStores(
	registry marketTypes.MarketRegistry,
	spotStore spotTypes.MarketStore,
) {
	registry.Register(spotStore)
}
