package testing

import (
	"github.com/wisp-trading/wisp/wisp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	wisp.Module,
)
