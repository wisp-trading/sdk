package adapters

import (
	"github.com/wisp-trading/sdk/pkg/adapters/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("adapters",
	fx.Provide(
		// Quiet default logger — no startup banner (CLI is a TUI).
		logging.NewDefaultZapLogger,
		logging.NewZapApplicationLogger,
		logging.NewZapTradingLogger,
	),
)
