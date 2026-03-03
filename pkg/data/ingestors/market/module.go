package market

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/perp"
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/spot"
	predictionRealtime "github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

var Module = fx.Module("market_ingestor",
	fx.Options(
		spot.Module,
		perp.Module,

		// Prediction realtime factory (from domain package)
		fx.Provide(
			fx.Annotate(
				predictionRealtime.NewFactory,
				fx.ResultTags(`group:"realtime_factories"`),
			),
		),

		fx.Provide(
			fx.Annotate(
				NewCoordinator,
				fx.ParamTags(`group:"batch_factories"`, `group:"realtime_factories"`),
			),
			newDataUpdateNotifier,
		),
	),
)

func newDataUpdateNotifier() ingestors.DataUpdateNotifier {
	return NewDataUpdateNotifier(100)
}
