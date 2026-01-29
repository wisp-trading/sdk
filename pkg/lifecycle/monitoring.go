package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/wisp-trading/sdk/pkg/monitoring"
	monitoringTypes "github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// initializeMonitoringServer creates and configures the monitoring server
// with remote shutdown capability via Unix socket HTTP endpoint
func (c *controller) initializeMonitoringServer(strategyName strategy.StrategyName) (monitoringTypes.Server, error) {
	return monitoring.NewServer(
		monitoringTypes.ServerConfig{
			InstanceID: string(strategyName),
		},
		c.viewRegistry,
		c.triggerShutdown, // Allow remote shutdown via HTTP endpoint
	)
}

// triggerShutdown handles remote shutdown requests from the monitoring server
// This is called when POST /shutdown is received on the Unix socket
//
// Critical for daemon operation: When strategies run as daemons without TTY,
// this provides the only way to gracefully stop them via the CLI/API
func (c *controller) triggerShutdown() {
	c.logger.Info("Remote shutdown triggered")

	// Create a background context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.Stop(ctx); err != nil {
		c.logger.Error("Error during shutdown: %v", err)
	}
}

// monitorHealth continuously monitors system health and reports aggregated errors
// Runs in a background goroutine and logs warnings when issues are detected
func (c *controller) monitorHealth(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			report := c.healthStore.GetSystemHealth()

			fmt.Println("Starting system health monitoring...")

			if report.HasErrors {
				c.logger.Warn("System health report:")

				// Log connector errors
				if len(report.ConnectorErrors.Errors) > 0 {
					c.logger.Warn("Connector errors:")
					for connector, err := range report.ConnectorErrors.Errors {
						c.logger.Warn("    - %s [%s]: %v", connector, err.State, err.Error)
					}
				}

				// Log data flow errors
				if len(report.DataFlowErrors.Errors) > 0 {
					c.logger.Warn("Data flow errors:")
					for connector, dataTypeErrors := range report.DataFlowErrors.Errors {
						for dataType, err := range dataTypeErrors {
							c.logger.Warn("    - %s:%s [%d errors]: %v", connector, dataType, err.ErrorCount, err.Error)
						}
					}
				}
			}
		}
	}
}
