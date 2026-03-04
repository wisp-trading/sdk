package market

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/spot"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	spotTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"go.uber.org/fx"
)

var Module = fx.Options(
	perp.Module,
	spot.Module,
	fx.Provide(
		NewMarketRegistry,
	),
	fx.Invoke(
		registerStores,
	),
)

// registerStores registers spot and perp stores with the registry.
// The prediction store is registered by pkg/markets/prediction/module.go.
func registerStores(
	registry marketTypes.MarketRegistry,
	spotStore spotTypes.MarketStore,
	perpStore perpTypes.MarketStore,
) {
	registry.Register(spotStore)
	registry.Register(perpStore)
}
