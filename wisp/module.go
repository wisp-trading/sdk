package wisp

import (
	packages "github.com/wisp-trading/sdk/pkg"
	"go.uber.org/fx"
)

// Module provides the Wisp SDK with all its services wired up via fx DI.
var Module = fx.Module("wisp",
	// Include all pkg modules
	packages.Module,

	// Provide the universe provider (caches trading universe)
	fx.Provide(NewUniverseProvider),

	// Provide the main Wisp context
	fx.Provide(NewWisp),
)
