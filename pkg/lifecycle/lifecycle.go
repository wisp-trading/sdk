package lifecycle

import (
	"context"
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

// controller implements the Controller interface for controlling SDK lifecycle
type controller struct {
	// Component references
	marketCoordinator   ingestors.MarketDataCoordinator
	positionCoordinator ingestors.PositionCoordinator
	connectorRegistry   registry.ConnectorRegistry
	logger              logging.ApplicationLogger

	// Lifecycle state
	state     lifecycleTypes.State
	stateMu   sync.RWMutex
	readyChan chan struct{}
	readyOnce sync.Once
}

// NewController creates a new SDK lifecycle controller
func NewController(
	marketCoordinator ingestors.MarketDataCoordinator,
	positionCoordinator ingestors.PositionCoordinator,
	connectorRegistry registry.ConnectorRegistry,
	logger logging.ApplicationLogger,
) lifecycleTypes.Controller {
	return &controller{
		marketCoordinator:   marketCoordinator,
		positionCoordinator: positionCoordinator,
		connectorRegistry:   connectorRegistry,
		logger:              logger,
		state:               lifecycleTypes.StateCreated,
		readyChan:           make(chan struct{}),
	}
}

// Start starts the SDK and all its components
func (c *controller) Start(ctx context.Context) error {
	c.stateMu.Lock()
	if c.state != lifecycleTypes.StateCreated && c.state != lifecycleTypes.StateStopped {
		currentState := c.state
		c.stateMu.Unlock()
		return fmt.Errorf("cannot start: current state is %v", currentState)
	}
	c.state = lifecycleTypes.StateStarting
	c.stateMu.Unlock()

	c.logger.Info("🚀 Initializing Kronos SDK...")

	// Wait for connectors to be ready
	if err := c.waitForConnectors(); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return err
	}

	// Start coordinators
	c.logger.Info("⚡ Starting data coordinators...")

	// Start position coordinator (if needed)
	if c.positionCoordinator != nil {
		c.logger.Info("  📊 Starting position tracking...")
		if err := c.positionCoordinator.Start(ctx); err != nil {
			c.stateMu.Lock()
			c.state = lifecycleTypes.StateCreated
			c.stateMu.Unlock()
			return fmt.Errorf("failed to start position coordinator: %w", err)
		}
		c.logger.Info("  ✓ Position tracking ready")
	}

	// Start market data ingestion
	c.logger.Info("  📈 Starting market data ingestion...")
	if err := c.marketCoordinator.StartDataCollection(ctx); err != nil {
		c.stateMu.Lock()
		c.state = lifecycleTypes.StateCreated
		c.stateMu.Unlock()
		return fmt.Errorf("failed to start market data collection: %w", err)
	}
	c.logger.Info("  ✓ Market data ingestion ready")

	// Mark as ready
	c.stateMu.Lock()
	c.state = lifecycleTypes.StateReady
	c.stateMu.Unlock()

	c.readyOnce.Do(func() {
		close(c.readyChan)
	})

	c.logger.Info("✅ Kronos SDK ready - strategies can now execute")
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

	c.logger.Info("Stopping Kronos SDK...")

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

	c.stateMu.Lock()
	c.state = lifecycleTypes.StateStopped
	c.stateMu.Unlock()

	c.logger.Info("Kronos SDK stopped")
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

// waitForConnectors waits for connectors to be marked ready
func (c *controller) waitForConnectors() error {
	c.logger.Info("🔌 Waiting for connectors to initialize...")

	readyConnectors := c.connectorRegistry.GetReadyConnectors()
	if len(readyConnectors) == 0 {
		c.logger.Warn("⚠️  No connectors marked as ready - SDK will start but data ingestion may fail")
		return nil
	}

	// Log each ready connector
	for _, conn := range readyConnectors {
		info := conn.GetConnectorInfo()
		c.logger.Info("  ✓ %s connected", info.Name)
	}

	c.logger.Info("✓ All %d connector(s) ready", len(readyConnectors))
	return nil
}
