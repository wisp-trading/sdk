package spot

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/spot/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/spot/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("spot_ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ResultTags(`name:"spot_batch_factory"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`name:"spot_realtime_factory"`),
		),
	),
)
