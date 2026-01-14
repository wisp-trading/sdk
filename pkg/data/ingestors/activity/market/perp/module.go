package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/perp/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/perp/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("perp_ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ResultTags(`name:"perp_batch_factory"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`name:"perp_realtime_factory"`),
		),
	),
)
