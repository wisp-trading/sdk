package signal

import (
	"go.uber.org/fx"
)

var Module = fx.Module("signal",
	fx.Provide(
		NewFactory,
	),
)
