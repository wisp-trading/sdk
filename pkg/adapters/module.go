package adapters

import (
	"github.com/backtesting-org/kronos-sdk/pkg/adapters/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("adaptors",
	fx.Provide(
		logging.NewZapApplicationLogger,
	),
)
