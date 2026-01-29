package spot

import (
	"github.com/wisp-trading/wisp/pkg/data/ingestors/market/spot/batch"
	"github.com/wisp-trading/wisp/pkg/data/ingestors/market/spot/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("spot_ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ParamTags(``, ``, `name:"spot_market_store"`),
			fx.ResultTags(`group:"batch_factories"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ParamTags(``, ``, `name:"spot_market_store"`),
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),
)
