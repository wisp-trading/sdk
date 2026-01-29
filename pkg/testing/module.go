package testing

import (
	"github.com/wisp-trading/sdk/wisp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	wisp.Module,
)
