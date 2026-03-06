package lifecycle

import (
	"context"
	"fmt"
	"sync"

	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	monitoringTypes "github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/health"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"go.uber.org/fx"
)

type controller struct {
	domains           []lifecycleTypes.DomainLifecycle
	connectorRegistry registry.ConnectorRegistry
	healthStore       health.HealthStore
	orchestrator      lifecycleTypes.Orchestrator
	logger            logging.ApplicationLogger
	viewRegistry      monitoringTypes.ViewRegistry
	monitoringServer  monitoringTypes.Server

	state     lifecycleTypes.State
	stateMu   sync.RWMutex
	readyChan chan struct{}
	readyOnce sync.Once
}

type controllerParams struct {
	fx.In

	Domains           []lifecycleTypes.DomainLifecycle `group:"domain_lifecycles"`
	ConnectorRegistry registry.ConnectorRegistry
	HealthStore       health.HealthStore
	Orchestrator      lifecycleTypes.Orchestrator
	Logger            logging.ApplicationLogger
	ViewRegistry      monitoringTypes.ViewRegistry
}

func NewController(p controllerParams) lifecycleTypes.Controller {
	return &controller{
		domains:           p.Domains,
		connectorRegistry: p.ConnectorRegistry,
		healthStore:       p.HealthStore,
		orchestrator:      p.Orchestrator,
		logger:            p.Logger,
		viewRegistry:      p.ViewRegistry,
		state:             lifecycleTypes.StateCreated,
		readyChan:         make(chan struct{}),
	}
}

func (c *controller) Start(ctx context.Context, strategyName strategy.StrategyName) error {
	c.stateMu.Lock()
	if c.state != lifecycleTypes.StateCreated && c.state != lifecycleTypes.StateStopped {
		current := c.state
		c.stateMu.Unlock()
		return fmt.Errorf("cannot start: current state is %v", current)
	}
	c.state = lifecycleTypes.StateStarting
	c.stateMu.Unlock()

	c.logger.Info("🚀 Initializing Wisp SDK...")

	if err := c.validateConnectorsReady(); err != nil {
		c.setState(lifecycleTypes.StateCreated)
		return err
	}

	// Start each domain in order — spot, perp, prediction each own their ingestors.
	for _, domain := range c.domains {
		c.logger.Info("  ⚡ Starting %s...", domain.Name())
		if err := domain.Start(ctx); err != nil {
			c.setState(lifecycleTypes.StateCreated)
			return fmt.Errorf("failed to start %s: %w", domain.Name(), err)
		}
		c.logger.Info("  ✓ %s ready", domain.Name())
	}

	c.logger.Info("  🎯 Starting strategy orchestrator...")
	if err := c.orchestrator.Start(ctx); err != nil {
		c.setState(lifecycleTypes.StateCreated)
		return fmt.Errorf("failed to start orchestrator: %w", err)
	}
	c.logger.Info("  ✓ Orchestrator ready")

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

	go c.monitorHealth(ctx)

	c.setState(lifecycleTypes.StateReady)
	c.readyOnce.Do(func() { close(c.readyChan) })

	c.logger.Info("✅ Wisp ready")
	return nil
}

func (c *controller) Stop(ctx context.Context) error {
	c.stateMu.Lock()
	if c.state == lifecycleTypes.StateStopped || c.state == lifecycleTypes.StateStopping {
		c.stateMu.Unlock()
		return nil
	}
	c.state = lifecycleTypes.StateStopping
	c.stateMu.Unlock()

	c.logger.Info("Stopping Wisp SDK...")

	if err := c.orchestrator.Stop(ctx); err != nil {
		c.logger.Error("Error stopping orchestrator: %v", err)
	}

	// Stop domains in reverse order.
	for i := len(c.domains) - 1; i >= 0; i-- {
		if err := c.domains[i].Stop(); err != nil {
			c.logger.Error("Error stopping %s: %v", c.domains[i].Name(), err)
		}
	}

	if c.monitoringServer != nil {
		if err := c.monitoringServer.Stop(ctx); err != nil {
			c.logger.Error("Error stopping monitoring server: %v", err)
		}
	}

	c.setState(lifecycleTypes.StateStopped)
	c.logger.Info("Wisp SDK stopped")
	return nil
}

func (c *controller) WaitUntilReady(ctx context.Context) error {
	select {
	case <-c.readyChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *controller) State() lifecycleTypes.State {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

func (c *controller) IsReady() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state == lifecycleTypes.StateReady
}

func (c *controller) setState(s lifecycleTypes.State) {
	c.stateMu.Lock()
	c.state = s
	c.stateMu.Unlock()
}

func (c *controller) validateConnectorsReady() error {
	ready := c.connectorRegistry.Filter(registry.NewFilter().ReadyOnly().Build())
	if len(ready) == 0 {
		return fmt.Errorf("no connectors marked as ready - initialize and mark connectors ready before calling Start()")
	}
	c.logger.Info("✓ Validated %d connector(s) ready", len(ready))
	return nil
}
