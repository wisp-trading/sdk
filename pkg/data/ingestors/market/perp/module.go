package perp

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/perp/batch"
	"github.com/wisp-trading/sdk/pkg/data/ingestors/market/perp/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("perp_ingestor",
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
