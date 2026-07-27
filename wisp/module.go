package wisp

import (
	sdkpkg "github.com/wisp-trading/sdk/pkg"
	"go.uber.org/fx"
)

// Module provides the Wisp SDK with all services wired via fx.
// Domain modules (spot/perp/options/prediction) each Provide their facade.New*;
// this module only composes the graph and NewWisp.
var Module = fx.Module("wisp",
	sdkpkg.Module,
	fx.Provide(NewWisp),
)
