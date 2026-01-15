package perp

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewStore,
			fx.ResultTags(`name:"perp_market_store"`),
		),
	),
)
