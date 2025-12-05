package runtime

import (
	"context"
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/plugin"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/runtime"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type rt struct {
	pluginManager     plugin.Manager
	connectorRegistry registry.ConnectorRegistry
	strategyRegistry  registry.StrategyRegistry
	controller        lifecycle.Controller
	logger            logging.ApplicationLogger
	loadedStrategy    strategy.Strategy
}

func NewRuntime(
	pluginManager plugin.Manager,
	connectorRegistry registry.ConnectorRegistry,
	strategyRegistry registry.StrategyRegistry,
	controller lifecycle.Controller,
	logger logging.ApplicationLogger,
) runtime.Runtime {
	return &rt{
		pluginManager:     pluginManager,
		connectorRegistry: connectorRegistry,
		strategyRegistry:  strategyRegistry,
		controller:        controller,
		logger:            logger,
	}
}

// Boot executes the complete startup sequence
func (r *rt) Boot(ctx context.Context, config runtime.BootConfig) error {
	r.logger.Info("🔧 Starting boot sequence...")

	// Step 1: Load strategy plugin
	r.logger.Info("Step 1/3: Loading strategy plugin...")
	strat, err := r.pluginManager.LoadStrategyPlugin(config.StrategyPath)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to load strategy plugin: %v", err))
		return fmt.Errorf("strategy load failed: %w", err)
	}

	r.loadedStrategy = strat
	r.logger.Info("✓ Strategy loaded: %s", strat.GetName())

	// Step 2: Register strategy
	r.logger.Info("Step 3/3: Starting SDK lifecycle...")
	r.strategyRegistry.RegisterStrategy(strat)

	// Step 3: Start SDK lifecycle
	if err := r.controller.Start(ctx); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to start controller: %v", err))
		return err
	}

	r.logger.Info("✅ Boot sequence complete - system ready")
	return nil
}

// Stop executes the graceful shutdown sequence
func (r *rt) Stop(ctx context.Context) error {
	r.logger.Info("🛑 Starting shutdown sequence...")

	// Stop SDK lifecycle
	if err := r.controller.Stop(ctx); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to stop controller: %v", err))
		return err
	}

	r.logger.Info("✅ Shutdown complete - system stopped")
	return nil
}
