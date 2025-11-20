package time

import (
	"go.uber.org/fx"
)

var Module = fx.Module("time",
	fx.Provide(
		NewTimeProvider,
	),
)
