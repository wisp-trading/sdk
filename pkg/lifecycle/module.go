package lifecycle

import (
	"go.uber.org/fx"
)

// Module provides the SDK lifecycle controller for dependency injection
var Module = fx.Module("lifecycle",
	fx.Provide(
		NewController,
		NewOrchestrator,
	),
)
