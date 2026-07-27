package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	runtimeTypes "github.com/wisp-trading/sdk/pkg/types/runtime"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// DefaultStopTimeout is how long Stop allows for domain/strategy teardown.
const DefaultStopTimeout = 30 * time.Second

type rt struct {
	connectorRegistry registry.ConnectorRegistry
	strategyRegistry  registry.StrategyRegistry
	configLoader      configTypes.StartupConfigLoader
	controller        lifecycle.Controller
	logger            logging.ApplicationLogger
	loadedStrategy    strategy.Strategy
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewRuntime(
	connectorRegistry registry.ConnectorRegistry,
	strategyRegistry registry.StrategyRegistry,
	configLoader configTypes.StartupConfigLoader,
	controller lifecycle.Controller,
	logger logging.ApplicationLogger,
) runtimeTypes.Runtime {
	return &rt{
		connectorRegistry: connectorRegistry,
		strategyRegistry:  strategyRegistry,
		configLoader:      configLoader,
		controller:        controller,
		logger:            logger,
	}
}

// StartStandalone runs a strategy in standalone mode (only packaging path).
// After a successful return, call Wait so /shutdown and OS signals share one stop path.
// settingsPath may be empty — uses ResolveSettingsPath (~/.wisp/connectors.yml).
func (r *rt) StartStandalone(
	strat strategy.Strategy,
	strategyDir string,
	settingsPath string,
) error {
	r.ctx, r.cancel = context.WithCancel(context.Background())

	cfg, err := r.configLoader.LoadForStrategy(strategyDir, settingsPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, err := r.initializeConnectors(cfg.ConnectorConfigs); err != nil {
		return err
	}

	return r.boot(r.ctx, cfg, strat)
}

// Wait blocks until OS signal or remote /shutdown, then stops the runtime.
func (r *rt) Wait() error {
	if r.ctx == nil {
		return fmt.Errorf("runtime not started")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case sig := <-sigChan:
		r.logger.Info("Received shutdown signal", "signal", sig.String())
	case <-r.controller.ShutdownRequested():
		r.logger.Info("Remote shutdown requested")
	case <-r.ctx.Done():
		r.logger.Info("Runtime context canceled")
	}

	return r.Stop()
}

// Stop gracefully shuts down using a fresh timeout context for cleanup,
// then cancels the long-lived root context.
func (r *rt) Stop() error {
	r.logger.Info("🛑 Stopping runtime...")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), DefaultStopTimeout)
	defer stopCancel()

	if err := r.controller.Stop(stopCtx); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to stop controller: %v", err))
		if r.cancel != nil {
			r.cancel()
		}
		return err
	}

	if r.cancel != nil {
		r.cancel()
	}

	r.logger.Info("✅ Runtime stopped")
	return nil
}

func (r *rt) initializeConnectors(connectors map[connector.ExchangeName]connector.Config) ([]connector.ExchangeName, error) {
	names := make([]connector.ExchangeName, 0, len(connectors))

	for name, cfg := range connectors {
		conn, exists := r.connectorRegistry.Connector(name)
		if !exists {
			r.logger.Warn(fmt.Sprintf("connector %s not registered", name))
			continue
		}

		if err := conn.Initialize(cfg); err != nil {
			return nil, fmt.Errorf("failed to initialize connector %s: %w", name, err)
		}

		if err := r.connectorRegistry.MarkReady(name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	r.logger.Info("Initialized connectors", "count", len(names))
	return names, nil
}

func (r *rt) boot(ctx context.Context, startupCfg *configTypes.StartupConfig, strat strategy.Strategy) error {
	if strat == nil {
		return fmt.Errorf("no strategy provided")
	}

	r.logger.Info("Booting...", "mode", runtimeTypes.BootModeStandalone)
	r.loadedStrategy = strat
	r.logger.Info("Strategy loaded", "name", strat.GetName())

	r.strategyRegistry.RegisterStrategy(strat)

	if err := r.controller.Start(ctx, strat.GetName(), startupCfg); err != nil {
		return fmt.Errorf("failed to start lifecycle: %w", err)
	}

	r.logger.Info("✅ Runtime started")
	return nil
}
