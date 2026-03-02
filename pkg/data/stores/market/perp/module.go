package perp

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewStore,
		),
	),
)
