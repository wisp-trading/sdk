package runtime

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// BootMode defines how strategies are loaded into the runtime
type BootMode string

const (
	BootModePlugin     BootMode = "plugin"
	BootModeStandalone BootMode = "standalone"
)

// BootConfig holds internal configuration for booting
type BootConfig struct {
	Mode           BootMode
	StrategyPath   string
	Strategy       strategy.Strategy
	ConnectorNames []connector.ExchangeName
}

// Runtime is the main entry point for running strategies.
//
// Blessed packaging path: StartStandalone + Wait (standalone binary with own main).
// Plugin mode (Start) is legacy and not recommended for new strategies.
type Runtime interface {
	// Start runs a strategy in plugin mode (legacy).
	// Prefer StartStandalone for new work. Loads config from configPath
	// (strategy dir) and wispPath (wisp.yml).
	Start(configPath string, wispPath string) error

	// StartStandalone runs a strategy in standalone mode (blessed path).
	// Use this from a strategy binary's main after fx wiring. After StartStandalone
	// returns successfully, call Wait so /shutdown and OS signals share one stop path.
	StartStandalone(strategy strategy.Strategy, configPath string, wispPath string) error

	// Wait blocks until an OS signal (SIGINT/SIGTERM) or remote /shutdown is
	// received, then performs a clean Stop and returns. Process hosts should
	// call Wait after a successful Start/StartStandalone so the process exits.
	Wait() error

	// Stop gracefully shuts down. Prefer Wait for the normal process contract;
	// Stop is still available for tests and custom hosts.
	Stop() error
}
