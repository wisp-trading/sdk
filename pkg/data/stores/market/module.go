package market

import (
	"github.com/wisp-trading/sdk/pkg/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/spot"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	predictionTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	spotTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"go.uber.org/fx"
)

var Module = fx.Options(
	perp.Module,
	spot.Module,
	prediction.Module,
	fx.Provide(
		NewMarketRegistry,
	),
	fx.Invoke(
		registerStores,
	),
)

// registerStores registers all market stores with the registry
func registerStores(
	registry marketTypes.MarketRegistry,
	spotStore spotTypes.MarketStore,
	perpStore perpTypes.MarketStore,
	predictionStore predictionTypes.MarketStore,
) {
	registry.Register(spotStore)
	registry.Register(perpStore)
	registry.Register(predictionStore)
}
