package runtime

import (
	"context"
	"fmt"

	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/plugin"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	runtimeTypes "github.com/wisp-trading/sdk/pkg/types/runtime"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type rt struct {
	pluginManager     plugin.Manager
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
	pluginManager plugin.Manager,
	connectorRegistry registry.ConnectorRegistry,
	strategyRegistry registry.StrategyRegistry,
	configLoader configTypes.StartupConfigLoader,
	controller lifecycle.Controller,
	logger logging.ApplicationLogger,
) runtimeTypes.Runtime {
	return &rt{
		pluginManager:     pluginManager,
		connectorRegistry: connectorRegistry,
		strategyRegistry:  strategyRegistry,
		configLoader:      configLoader,
		controller:        controller,
		logger:            logger,
	}
}

// Start runs a strategy in plugin mode.
func (r *rt) Start(configPath string, wispPath string) error {
	r.ctx, r.cancel = context.WithCancel(context.Background())

	cfg, err := r.configLoader.LoadForStrategy(configPath, wispPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, err := r.initializeConnectors(cfg.ConnectorConfigs); err != nil {
		return err
	}

	return r.boot(r.ctx, cfg, runtimeTypes.BootConfig{
		Mode:         runtimeTypes.BootModePlugin,
		StrategyPath: cfg.PluginPath,
	})
}

// StartStandalone runs a strategy in standalone mode (debuggable).
func (r *rt) StartStandalone(
	strat strategy.Strategy,
	configPath string,
	wispPath string,
) error {
	r.ctx, r.cancel = context.WithCancel(context.Background())

	cfg, err := r.configLoader.LoadForStrategy(configPath, wispPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, err := r.initializeConnectors(cfg.ConnectorConfigs); err != nil {
		return err
	}

	return r.boot(r.ctx, cfg, runtimeTypes.BootConfig{
		Mode:     runtimeTypes.BootModeStandalone,
		Strategy: strat,
	})
}

// Stop gracefully shuts down.
func (r *rt) Stop() error {
	r.logger.Info("🛑 Stopping runtime...")

	if r.cancel != nil {
		r.cancel()
	}

	if err := r.controller.Stop(r.ctx); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to stop controller: %v", err))
		return err
	}

	r.logger.Info("✅ Runtime stopped")
	return nil
}

// initializeConnectors initializes all connectors and marks them ready.
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

// boot loads the strategy and hands off to the lifecycle controller.
func (r *rt) boot(ctx context.Context, startupCfg *configTypes.StartupConfig, cfg runtimeTypes.BootConfig) error {
	r.logger.Info("Booting...", "mode", cfg.Mode)

	var (
		strat strategy.Strategy
		err   error
	)

	switch cfg.Mode {
	case runtimeTypes.BootModeStandalone:
		if cfg.Strategy == nil {
			return fmt.Errorf("no strategy provided in standalone mode")
		}
		strat = cfg.Strategy

	default:
		r.logger.Info("Loading strategy plugin...", "path", cfg.StrategyPath)
		strat, err = r.pluginManager.LoadStrategyPlugin(cfg.StrategyPath)
		if err != nil {
			return fmt.Errorf("failed to load plugin: %w", err)
		}
	}

	r.loadedStrategy = strat
	r.logger.Info("Strategy loaded", "name", strat.GetName())

	r.strategyRegistry.RegisterStrategy(strat)

	if err := r.controller.Start(ctx, strat.GetName(), startupCfg); err != nil {
		return fmt.Errorf("failed to start lifecycle: %w", err)
	}

	r.logger.Info("✅ Runtime started")
	return nil
}
