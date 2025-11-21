package lifecycle

import (
	"context"

	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"go.uber.org/fx"
)

// Module provides the SDK lifecycle controller for dependency injection
var Module = fx.Module("lifecycle",
	fx.Provide(NewController),

	// Register lifecycle hooks to automatically start/stop the controller
	fx.Invoke(registerLifecycleHooks),
)

func registerLifecycleHooks(lc fx.Lifecycle, controller lifecycleTypes.Controller, logger logging.ApplicationLogger) {
	// Create a long-lived context that will be used for the entire lifecycle
	// This context is separate from the fx OnStart context which is short-lived
	lifecycleCtx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start the controller with the long-lived context
			logger.Info("Initializing Kronos SDK lifecycle controller")
			if err := controller.Start(lifecycleCtx); err != nil {
				cancel()
				logger.Error("Failed to start lifecycle controller: %v", err)
				return err
			}
			logger.Info("Lifecycle controller started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down lifecycle controller")
			// Cancel the long-lived context and stop the controller
			cancel()
			if err := controller.Stop(ctx); err != nil {
				logger.Error("Error stopping lifecycle controller: %v", err)
				return err
			}
			logger.Info("Lifecycle controller stopped successfully")
			return nil
		},
	})
}
