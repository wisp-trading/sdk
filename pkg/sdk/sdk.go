package sdk

import (
	"context"
	"fmt"

	kronosPackage "github.com/backtesting-org/kronos-sdk/kronos"
	kronosTypes "github.com/backtesting-org/kronos-sdk/pkg/types/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/runtime"
	"go.uber.org/fx"
)

// Initialize creates and starts the Kronos SDK for standalone use.
// This is a convenience function for standalone binaries that want to run
// strategies directly without the plugin system.
//
// Returns:
//   - runtime.Runtime: The runtime instance for booting strategies
//   - kronosTypes.Kronos: The Kronos context for strategy creation
//   - func(): Cleanup function to stop the SDK
//   - error: Any initialization error
//
// Example usage:
//
//	import "github.com/backtesting-org/kronos-sdk/pkg/sdk"
//
//	ctx := context.Background()
//	rt, k, cleanup, err := sdk.Initialize(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to initialize: %v", err)
//	}
//	defer cleanup()
//
//	strategy := mystrategy.New(k)
//	err = rt.Boot(ctx, runtime.BootConfig{
//	    Mode:     runtime.BootModeStandalone,
//	    Strategy: strategy,
//	})
func Initialize(ctx context.Context) (runtime.Runtime, kronosTypes.Kronos, func(), error) {
	var rt runtime.Runtime
	var k kronosTypes.Kronos

	app := fx.New(
		kronosPackage.Module,
		fx.Populate(&rt, &k),
		fx.NopLogger,
	)

	if err := app.Start(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start Kronos SDK: %w", err)
	}

	cleanup := func() {
		if err := app.Stop(ctx); err != nil {
			// Log but don't panic on cleanup errors
			fmt.Printf("Warning: error during SDK cleanup: %v\n", err)
		}
	}

	return rt, k, cleanup, nil
}
