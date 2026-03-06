package runtime

import (
	"context"
	"fmt"

	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	configTypes "github.com/wisp-trading/sdk/pkg/types/config"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/plugin"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/runtime"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type rt struct {
	pluginManager     plugin.Manager
	connectorRegistry registry.ConnectorRegistry
	marketWatchlist   types.MarketWatchlist
	perpWatchlist     perpTypes.PerpWatchlist
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
	marketWatchlist types.MarketWatchlist,
	perpWatchlist perpTypes.PerpWatchlist,
	strategyRegistry registry.StrategyRegistry,
	configLoader configTypes.StartupConfigLoader,
	controller lifecycle.Controller,
	logger logging.ApplicationLogger,
) runtime.Runtime {
	return &rt{
		pluginManager:     pluginManager,
		connectorRegistry: connectorRegistry,
		marketWatchlist:   marketWatchlist,
		perpWatchlist:     perpWatchlist,
		strategyRegistry:  strategyRegistry,
		configLoader:      configLoader,
		controller:        controller,
		logger:            logger,
	}
}

// Start runs a strategy in plugin mode
func (r *rt) Start(configPath string, wispPath string) error {
	r.ctx, r.cancel = context.WithCancel(context.Background())

	cfg, err := r.configLoader.LoadForStrategy(configPath, wispPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	connectorNames, err := r.initializeConnectors(cfg.ConnectorConfigs)
	if err != nil {
		return err
	}

	r.registerAssets(cfg)

	return r.boot(r.ctx, runtime.BootConfig{
		Mode:           runtime.BootModePlugin,
		StrategyPath:   cfg.PluginPath,
		ConnectorNames: connectorNames,
	})
}

// StartStandalone runs a strategy in standalone mode (debuggable)
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

	connectorNames, err := r.initializeConnectors(cfg.ConnectorConfigs)
	if err != nil {
		return err
	}

	r.registerAssets(cfg)

	return r.boot(r.ctx, runtime.BootConfig{
		Mode:           runtime.BootModeStandalone,
		Strategy:       strat,
		ConnectorNames: connectorNames,
	})
}

// Stop gracefully shuts down
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

// initializeConnectors initializes all connectors and returns their names
func (r *rt) initializeConnectors(connectors map[connector.ExchangeName]connector.Config) ([]connector.ExchangeName, error) {
	names := make([]connector.ExchangeName, 0, len(connectors))

	for name, cfg := range connectors {
		conn, exists := r.connectorRegistry.Connector(name)
		if !exists {
			r.logger.Warn(fmt.Sprintf("connector %s not registered", name))
			continue
		}

		if err := conn.Initialize(cfg); err != nil {
			r.logger.Error(fmt.Sprintf("connector %s init failed: %s", name, err.Error()))
			return nil, fmt.Errorf("failed to initialize connector %s: %w", name, err)
		}

		names = append(names, name)
		if err := r.connectorRegistry.MarkReady(name); err != nil {
			return nil, err
		}
	}

	r.logger.Info("Initialized connectors", "count", len(names))
	return names, nil
}

// registerAssets routes config assets to the correct domain watchlist based on
// the registered connector type — no market_type annotation needed in config.
// Must be called after initializeConnectors so connector types are known.
func (r *rt) registerAssets(cfg *configTypes.StartupConfig) {
	spotCount, perpCount, unknownCount := 0, 0, 0

	for exchange, pairs := range cfg.Assets {
		marketType, ok := r.connectorRegistry.ConnectorType(exchange)
		if !ok {
			r.logger.Warn("No connector registered for exchange %s — skipping asset registration", exchange)
			unknownCount += len(pairs)
			continue
		}

		for _, pair := range pairs {
			switch marketType {
			case connector.MarketTypePerp:
				r.perpWatchlist.RequirePair(exchange, pair)
				perpCount++
			case connector.MarketTypeSpot:
				r.marketWatchlist.RequirePair(exchange, pair)
				spotCount++
			default:
				// Prediction markets don't use a pair watchlist
				r.logger.Debug("Skipping pair %s for %s connector on %s",
					pair.Symbol(), marketType, exchange)
			}
		}
	}

	r.logger.Info("Registered assets from config",
		"spot", spotCount, "perp", perpCount, "skipped", unknownCount)
}

// boot executes the core boot sequence
func (r *rt) boot(ctx context.Context, cfg runtime.BootConfig) error {
	r.logger.Info("🔧 Booting...", "mode", cfg.Mode)

	var strat strategy.Strategy
	var err error

	switch cfg.Mode {
	case runtime.BootModeStandalone:
		r.logger.Info("Using provided strategy instance")
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

	if err := r.controller.Start(ctx, strat.GetName()); err != nil {
		return fmt.Errorf("failed to start lifecycle: %w", err)
	}

	r.logger.Info("✅ Runtime started")
	return nil
}
