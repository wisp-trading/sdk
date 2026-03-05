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

// Runtime is the main entry point for running strategies
type Runtime interface {
	// Start runs a strategy in plugin mode
	// Loads config from configPath (strategy dir) and wispPath (wisp.yml)
	Start(configPath string, wispPath string) error

	// StartStandalone runs a strategy in standalone mode (debuggable)
	StartStandalone(strategy strategy.Strategy, configPath string, wispPath string) error

	// Stop gracefully shuts down
	Stop() error
}
