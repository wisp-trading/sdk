package registry

import (
	"go.uber.org/fx"
)

// Module provides registry implementations
var Module = fx.Module("registry",
	fx.Provide(
		NewConnectorRegistry,
		NewStrategyRegistry,
		NewHookRegistry,
	),
)
