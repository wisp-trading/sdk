package runtime

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
)

// BootConfig holds configuration for booting the runtime
type BootConfig struct {
	StrategyPath string

	// Connectors to mark as ready - must be initialized by caller before Boot
	ConnectorNames []connector.ExchangeName
}

// Runtime orchestrates the complete boot sequence:
// 1. Load strategy plugin
// 2. Verify connectors are registered
// 3. Start the SDK lifecycle
type Runtime interface {
	Boot(ctx context.Context, config BootConfig) error
	Stop(ctx context.Context) error
}
