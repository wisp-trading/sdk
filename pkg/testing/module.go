package testing

import (
	packages "github.com/backtesting-org/kronos-sdk/pkg"
	"go.uber.org/fx"
)

var Module = fx.Options(
	packages.Module,
)
