package market

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/perp"
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/prediction"
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/spot"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

var Module = fx.Module("market_ingestor",
	fx.Options(
		spot.Module,
		perp.Module,
		prediction.Module,

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
