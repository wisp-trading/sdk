package perp

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/perp/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/perp/realtime"
	"go.uber.org/fx"
)

var Module = fx.Module("perp_ingestor",
	fx.Provide(
		fx.Annotate(
			batch.NewFactory,
			fx.ParamTags(``, ``, `name:"perp_market_store"`),
			fx.ResultTags(`group:"batch_factories"`),
		),
		fx.Annotate(
			realtime.NewFactory,
			fx.ParamTags(``, ``, `name:"perp_market_store"`),
			fx.ResultTags(`group:"realtime_factories"`),
		),
	),
)
