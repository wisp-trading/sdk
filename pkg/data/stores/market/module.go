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
	fx.Invoke(fx.Annotate(
		registerStores,
		fx.ParamTags(``, `name:"spot_market_store"`, `name:"perp_market_store"`),
	)),
)

// registerStores registers all market stores with the registry
func registerStores(
	registry marketTypes.MarketRegistry,
	spotStore spotTypes.MarketStore,
	perpStore perpTypes.MarketStore,
) {
	registry.Register(spotStore)
	registry.Register(perpStore)
}
