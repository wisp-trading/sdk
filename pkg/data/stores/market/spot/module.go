package spot

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewStore,
			fx.ResultTags(`name:"spot_market_store"`),
		),
	),
)
