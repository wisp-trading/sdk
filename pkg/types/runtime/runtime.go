package runtime

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// BootMode defines how strategies are loaded into the runtime
type BootMode string

const (
	// BootModeStandalone is the only supported packaging path: strategy owns main.
	BootModeStandalone BootMode = "standalone"
)

// BootConfig holds internal configuration for booting
type BootConfig struct {
	Mode           BootMode
	Strategy       strategy.Strategy
	ConnectorNames []connector.ExchangeName
}

// Runtime is the main entry point for running strategies.
//
// Packaging path: StartStandalone + Wait (standalone binary with own main).
// Plugin / .so loading has been removed.
type Runtime interface {
	// StartStandalone runs a strategy in standalone mode.
	// Use from a strategy binary's main after fx wiring. After success, call Wait
	// so /shutdown and OS signals share one stop path.
	// settingsPath may be empty to use ~/.wisp/connectors.yml (or migration paths).
	StartStandalone(strategy strategy.Strategy, strategyDir string, settingsPath string) error

	// Wait blocks until an OS signal (SIGINT/SIGTERM) or remote /shutdown is
	// received, then performs a clean Stop and returns.
	Wait() error

	// Stop gracefully shuts down. Prefer Wait for the normal process contract.
	Stop() error
}
