package prediction

import (
	realtime "github.com/wisp-trading/sdk/pkg/data/ingestors/market/prediction/real_time"
	"go.uber.org/fx"
)

var Module = fx.Module("prediction_ingestor",
	fx.Provide(
		//fx.Annotate(
		//	batch.NewFactory,
		//	fx.ResultTags(`group:"batch_factories"`),
		//),
		fx.Annotate(
			realtime.NewFactory,
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),
)
