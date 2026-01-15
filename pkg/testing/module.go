package testing

import (
	"github.com/backtesting-org/kronos-sdk/kronos"
	"go.uber.org/fx"
)

var Module = fx.Options(
	kronos.Module,
)
