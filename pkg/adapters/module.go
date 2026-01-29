package adapters

import (
	"github.com/wisp-trading/wisp/pkg/adapters/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("adapters",
	fx.Provide(
		// Provide default zap logger if none exists
		fx.Annotate(
			logging.NewDefaultZapLogger,
			fx.OnStart(func(logger *zap.Logger) error {
				logger.Info("Zap logger initialized")
				return nil
			}),
		),
		logging.NewZapApplicationLogger,
		logging.NewZapTradingLogger,
	),
)
