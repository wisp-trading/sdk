package wisp

import (
	packages "github.com/wisp-trading/sdk/pkg"
	"github.com/wisp-trading/sdk/pkg/markets/perp/perp"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/predict"
	"go.uber.org/fx"
)

// Module provides the Wisp SDK with all its services wired up via fx DI.
var Module = fx.Module("wisp",
	// Include all pkg modules
	packages.Module,

	// Provide the universe provider (caches trading universe)
	fx.Provide(NewUniverseProvider),

	// Provide domain context objects exposed on wisp
	fx.Provide(
		perp.NewPerp,
		predict.NewPredict,
	),

	// Provide the main Wisp context
	fx.Provide(NewWisp),
)
