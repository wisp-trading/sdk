package market

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/spot"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors"
	predictionRealtime "github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("market_ingestor",
	fx.Options(
		spot.Module,

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
