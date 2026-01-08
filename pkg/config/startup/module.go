package startup

import "go.uber.org/fx"

var Module = fx.Module("startup-config",
	fx.Provide(NewStartupConfigLoader),
)
