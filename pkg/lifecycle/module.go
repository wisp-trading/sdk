package lifecycle

import (
	"context"

	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"go.uber.org/fx"
)

// Module provides the SDK lifecycle controller for dependency injection
var Module = fx.Module("lifecycle",
	fx.Provide(NewController),

	// Register lifecycle but DON'T auto-start
	fx.Invoke(registerLifecycleHooks),
)

func registerLifecycleHooks(lc fx.Lifecycle, controller lifecycleTypes.Controller) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Controller is created but NOT started
			// Orchestrator/application must call controller.Start() explicitly
			// This hook just ensures the controller is available
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Always cleanup on shutdown
			return controller.Stop(ctx)
		},
	})
}
