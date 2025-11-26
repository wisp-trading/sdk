package lifecycle

import (
	"context"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type orchestrator struct {
	executor         execution.Executor
	strategyRegistry registry.StrategyRegistry
	logger           logging.ApplicationLogger
	timeProvider     temporal.TimeProvider
	notifier         ingestors.DataUpdateNotifier

	// Execution control
	ctx    context.Context
	cancel context.CancelFunc

	// Tick timer for data-driven execution
	tickTimer *TickTimer

	// Prevent concurrent executions of same strategy
	strategyMutexes map[strategy.StrategyName]*sync.Mutex
	mutexMapLock    sync.RWMutex
}

// NewOrchestrator creates a new strategy orchestrator
func NewOrchestrator(
	executor execution.Executor,
	strategyRegistry registry.StrategyRegistry,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	notifier ingestors.DataUpdateNotifier,
) lifecycleTypes.Orchestrator {
	tickTimer := NewTickTimer(
		5,                    // Execute after 5 data updates
		5*time.Second,        // Fallback every 5 seconds
		100*time.Millisecond, // Minimum 100ms between executions
		timeProvider,
	)

	orch := &orchestrator{
		executor:         executor,
		strategyRegistry: strategyRegistry,
		logger:           logger,
		timeProvider:     timeProvider,
		tickTimer:        tickTimer,
		strategyMutexes:  make(map[strategy.StrategyName]*sync.Mutex),
		notifier:         notifier,
	}

	return orch
}

// listenForDataUpdates forwards data update notifications to the tick timer
func (o *orchestrator) listenForDataUpdates() {
	for {
		select {
		case <-o.ctx.Done():
			// Orchestrator stopped, exit goroutine
			return
		case _, ok := <-o.notifier.Updates():
			if !ok {
				// Channel closed, exit goroutine
				return
			}
			// Forward notification to tick timer
			o.tickTimer.NotifyDataUpdate()
		}
	}
}

// Start begins orchestration
func (o *orchestrator) Start(ctx context.Context) error {
	if o.cancel != nil {
		o.logger.Warn("Orchestrator already started")
		return nil
	}

	o.ctx, o.cancel = context.WithCancel(ctx)
	o.logger.Info("🎯 Starting strategy orchestrator")

	go o.listenForDataUpdates()
	go o.executionLoop()

	o.logger.Info("✅ Strategy orchestrator started")
	return nil
}

// Stop gracefully stops orchestration
func (o *orchestrator) Stop(ctx context.Context) error {
	if o.cancel == nil {
		return nil
	}

	o.logger.Info("🛑 Stopping strategy orchestrator")

	// Stop the tick timer first to clean up its goroutine
	o.tickTimer.Stop()

	// Then cancel the orchestrator's context
	o.cancel()
	o.cancel = nil

	o.logger.Info("✅ Strategy orchestrator stopped")
	return nil
}

// NotifyDataUpdate triggers strategy execution on new market data
func (o *orchestrator) NotifyDataUpdate() {
	o.tickTimer.NotifyDataUpdate()
}

// GetStrategies returns all registered strategies
func (o *orchestrator) GetStrategies() []strategy.Strategy {
	return o.strategyRegistry.GetAllStrategies()
}

// executionLoop handles strategy execution triggers
func (o *orchestrator) executionLoop() {
	for {
		select {
		case <-o.ctx.Done():
			return
		case <-o.tickTimer.TickChannel():
			o.executeEnabledStrategies()
		}
	}
}

// executeEnabledStrategies runs all registered strategies
func (o *orchestrator) executeEnabledStrategies() {
	strategies := o.strategyRegistry.GetAllStrategies()

	if len(strategies) == 0 {
		return
	}

	// Execute each strategy concurrently
	var wg sync.WaitGroup
	for _, strat := range strategies {
		wg.Add(1)
		go func(s strategy.Strategy) {
			defer wg.Done()
			o.executeStrategy(s)
		}(strat)
	}
	wg.Wait()
}

// getStrategyMutex returns a mutex for the given strategy, creating one if needed
func (o *orchestrator) getStrategyMutex(strategyName strategy.StrategyName) *sync.Mutex {
	o.mutexMapLock.RLock()
	mutex, exists := o.strategyMutexes[strategyName]
	o.mutexMapLock.RUnlock()

	if exists {
		return mutex
	}

	o.mutexMapLock.Lock()
	defer o.mutexMapLock.Unlock()

	// Double-check after acquiring write lock
	if mutex, exists := o.strategyMutexes[strategyName]; exists {
		return mutex
	}

	mutex = &sync.Mutex{}
	o.strategyMutexes[strategyName] = mutex
	return mutex
}

// executeStrategy runs a single strategy
func (o *orchestrator) executeStrategy(strat strategy.Strategy) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("Strategy %s panicked: %v", strat.GetName(), r)
		}
	}()

	// Prevent concurrent executions of the same strategy
	strategyMutex := o.getStrategyMutex(strat.GetName())
	strategyMutex.Lock()
	defer strategyMutex.Unlock()

	startTime := o.timeProvider.Now()
	signals, err := strat.GetSignals()
	duration := o.timeProvider.Since(startTime)

	if err != nil {
		o.logger.Error("Strategy %s failed: %v (duration: %v)", strat.GetName(), err, duration)
		return
	}

	if len(signals) == 0 {
		o.logger.Debug("Strategy %s generated no signals (duration: %v)", strat.GetName(), duration)
		return
	}

	o.logger.Info("✅ Strategy %s generated %d signal(s) (duration: %v)", strat.GetName(), len(signals), duration)

	// Execute signals
	for _, signal := range signals {
		if err := o.executor.ExecuteSignal(signal); err != nil {
			o.logger.Error("Signal execution failed for %s: %v", strat.GetName(), err)
		}
	}
}
