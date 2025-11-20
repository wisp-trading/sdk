package registry

import "go.uber.org/fx"

var Module = fx.Provide(
	NewConnectorRegistry(),
)
