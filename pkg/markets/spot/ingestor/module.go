package ingestor

import (
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/spot/ingestor/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("spot_ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ResultTags(`group:"batch_factories"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),
)
