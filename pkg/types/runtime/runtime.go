package runtime

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// BootMode defines how strategies are loaded into the runtime
type BootMode string

const (
	// BootModePlugin loads strategies from compiled .so plugin files
	BootModePlugin BootMode = "plugin"
	// BootModeStandalone registers strategies directly (for debugging and single-binary deployments)
	BootModeStandalone BootMode = "standalone"
)

// BootConfig holds configuration for booting the runtime
type BootConfig struct {
	// Mode determines how the strategy is loaded
	Mode BootMode

	// StrategyPath is the path to the .so plugin file (used in BootModePlugin)
	StrategyPath string

	// Strategy is the directly provided strategy instance (used in BootModeStandalone)
	Strategy strategy.Strategy

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
