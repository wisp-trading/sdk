package plugin

import (
	"go.uber.org/fx"
)

// Module provides the plugin manager for dependency injection
var Module = fx.Module("plugin",
	fx.Provide(
		NewManager,
	),
)
