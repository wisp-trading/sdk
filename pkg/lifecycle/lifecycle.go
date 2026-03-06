package lifecycle

import (
	"context"
	"fmt"
	"sync"

	ingestors2 "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	monitoringTypes "github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/health"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// controller implements the Controller interface for controlling SDK lifecycle
type controller struct {
	// Component references
	marketCoordinator   ingestors2.MarketDataCoordinator
	positionCoordinator ingestors2.PositionCoordinator
	connectorRegistry   registry.ConnectorRegistry
	healthStore         health.HealthStore
	orchestrator        lifecycleTypes.Orchestrator
	logger              logging.ApplicationLogger
	viewRegistry        monitoringTypes.ViewRegistry
	monitoringServer    monitoringTypes.Server

	// Lifecycle state
	state     lifecycleTypes.State
	stateMu   sync.RWMutex
	readyChan chan struct{}
	readyOnce sync.Once
}

// NewController creates a new SDK lifecycle controller
func NewController(
	marketCoordinator ingestors2.MarketDataCoordinator,
	positionCoordinator ingestors2.PositionCoordinator,
	connectorRegistry registry.ConnectorRegistry,
	healthStore health.HealthStore,
	orchestrator lifecycleTypes.Orchestrator,
	logger logging.ApplicationLogger,
	viewRegistry monitoringTypes.ViewRegistry,
) lifecycleTypes.Controller {
	return &controller{
		marketCoordinator:   marketCoordinator,
		positionCoordinator: positionCoordinator,
		connectorRegistry:   connectorRegistry,
		healthStore:         healthStore,
		orchestrator:        orchestrator,
		logger:              logger,
		viewRegistry:        viewRegistry,
		state:               lifecycleTypes.StateCreated,
		readyChan:           make(chan struct{}),
	}
}

// Start starts the SDK and all its components
func (c *controller) Start(ctx context.Context, strategyName strategy.StrategyName) error {
	c.stateMu.Lock()
	if c.state != lifecycleTypes.StateCreated && c.state != lifecycleTypes.StateStopped {
		currentState := c.state
		c.stateMu.Unlock()
		return fmt.Errorf("cannot start: current state is %v", currentState)
	}
	c.state = lifecycleTypes.StateStarting
	c.stateMu.Unlock()

	c.logger.Info("🚀 Initializing Wisp SDK...")

	// Validate connectors are ready (must be initialized before Start())
	if err := c.validateConnectorsReady(); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return err
	}

	// Start coordinators
	c.logger.Info("⚡ Starting data coordinators...")

	// Start position coordinator (if needed)
	c.logger.Info("  📊 Starting position tracking...")
	if err := c.positionCoordinator.Start(ctx); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return fmt.Errorf("failed to start position coordinator: %w", err)
	}
	c.logger.Info("  ✓ Position tracking ready")

	// Start market data ingestion
	c.logger.Info("  📈 Starting market data ingestion...")
	if err := c.marketCoordinator.StartDataCollection(ctx); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return fmt.Errorf("failed to start market data collection: %w", err)
	}

	c.logger.Info("  ✓ Market data ingestion ready")

	// Start orchestrator
	c.logger.Info("   Starting strategy orchestrator...")
	if err := c.orchestrator.Start(ctx); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return fmt.Errorf("failed to start orchestrator: %w", err)
	}
	c.logger.Info("  ✓ Strategy orchestrator ready")

	// Start monitoring server
	c.logger.Info("Starting monitoring server...")
	monitoringServer, err := c.initializeMonitoringServer(strategyName)
	if err != nil {
		c.logger.Warn("Failed to initialize monitoring server: %v (continuing without monitoring)", err)
	} else {
		c.monitoringServer = monitoringServer
		go func() {
			if err := c.monitoringServer.Start(); err != nil {
				c.logger.Error("Monitoring server error: %v", err)
			}
		}()
		c.logger.Info("Monitoring server ready on %s", c.monitoringServer.SocketPath())
	}

	// Start runtime health monitoring
	go c.monitorHealth(ctx)

	// Mark as ready
	c.stateMu.Lock()
	c.state = lifecycleTypes.StateReady
	c.stateMu.Unlock()

	c.readyOnce.Do(func() {
		close(c.readyChan)
	})

	c.logger.Info("Wisp ready")
	return nil
}

// Stop gracefully shuts down the SDK
func (c *controller) Stop(ctx context.Context) error {
	c.stateMu.Lock()
	if c.state == lifecycleTypes.StateStopped || c.state == lifecycleTypes.StateStopping {
		c.stateMu.Unlock()
		return nil
	}
	c.state = lifecycleTypes.StateStopping
	c.stateMu.Unlock()

	c.logger.Info("Stopping Wisp SDK...")

	// Stop orchestrator first (stop generating new signals)
	if err := c.orchestrator.Stop(ctx); err != nil {
		c.logger.Error("Error stopping orchestrator: %v", err)
	}

	// Stop market data ingestion
	if err := c.marketCoordinator.StopDataCollection(); err != nil {
		c.logger.Error("Error stopping market data collection: %v", err)
	}

	// Stop position coordinator
	if c.positionCoordinator != nil {
		if err := c.positionCoordinator.Stop(); err != nil {
			c.logger.Error("Error stopping position coordinator: %v", err)
		}
	}

	// Stop monitoring server
	if c.monitoringServer != nil {
		c.logger.Info("Stopping monitoring server...")
		if err := c.monitoringServer.Stop(ctx); err != nil {
			c.logger.Error("Error stopping monitoring server: %v", err)
		}
	}

	c.stateMu.Lock()
	c.state = lifecycleTypes.StateStopped
	c.stateMu.Unlock()

	c.logger.Info("Wisp SDK stopped")
	return nil
}

// WaitUntilReady blocks until the SDK is ready or context is cancelled
func (c *controller) WaitUntilReady(ctx context.Context) error {
	select {
	case <-c.readyChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// State returns the current lifecycle state
func (c *controller) State() lifecycleTypes.State {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

// IsReady returns true if the SDK is ready
func (c *controller) IsReady() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state == lifecycleTypes.StateReady
}

// validateConnectorsReady validates that at least one connector is ready.
// Applications must initialize and mark connectors ready BEFORE calling Start().
// This method does not wait - it only verifies the precondition.
func (c *controller) validateConnectorsReady() error {
	readyConnectors := c.connectorRegistry.Filter(registry.NewFilter().ReadyOnly().Build())
	if len(readyConnectors) == 0 {
		return fmt.Errorf("no connectors marked as ready - initialize and mark connectors ready before calling Start()")
	}

	c.logger.Info("✓ Validated %d connector(s) ready", len(readyConnectors))

	return nil
}
