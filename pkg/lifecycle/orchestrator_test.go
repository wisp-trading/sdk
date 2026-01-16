package lifecycle_test

import (
	"context"
	"errors"
	"time"

	mockStrategy "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Orchestrator", func() {
	var (
		app              *fxtest.App
		orchestrator     lifecycleTypes.Orchestrator
		strategyRegistry registry.StrategyRegistry
		executor         execution.Executor
		notifier         ingestors.DataUpdateNotifier
		ctx              strategy.StrategyContext
	)

	BeforeEach(func() {
		ctx = strategy.NewStrategyContext(context.Background(), strategy.StrategyName("TestStrategy"))

		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&orchestrator,
				&strategyRegistry,
				&executor,
				&notifier,
			),
			fx.NopLogger,
		)
		app.RequireStart()
	})

	AfterEach(func() {
		if orchestrator != nil {
			_ = orchestrator.Stop(ctx)
		}
		app.RequireStop()
	})

	Describe("Starting and Stopping", func() {
		It("should start successfully", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not allow starting twice", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Start(ctx)
			// Should warn but not error
			Expect(err).ToNot(HaveOccurred())
		})

		It("should stop successfully", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Stop(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle multiple stop calls gracefully", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Stop(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Stop(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Data Update Notification Flow", func() {
		It("should forward data updates to tick timer", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Trigger data update through notifier
			notifier.Notify()

			// Allow time for processing
			time.Sleep(50 * time.Millisecond)
		})
	})

	Describe("Strategy Execution", func() {
		It("should execute single strategy", func() {
			// Setup mock strategy
			mockStrat := mockStrategy.NewStrategy(GinkgoT())
			mockStrat.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).Return([]*strategy.Signal{}, nil).Once()

			strategyRegistry.RegisterStrategy(mockStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Wait for MinExecutionInterval to pass (100ms default)
			time.Sleep(110 * time.Millisecond)

			// Trigger execution - need to reach threshold (default 5)
			for i := 0; i < 5; i++ {
				orchestrator.NotifyDataUpdate()
			}

			// Allow time for async execution
			time.Sleep(50 * time.Millisecond)
		})

		It("should execute multiple strategies concurrently", func() {
			// Setup mock strategies
			mockStrat1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrat2 := mockStrategy.NewStrategy(GinkgoT())

			mockStrat1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat1.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).Return([]*strategy.Signal{}, nil).Once()
			mockStrat2.EXPECT().GetName().Return(strategy.StrategyName("Strategy2")).Maybe()
			mockStrat2.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).Return([]*strategy.Signal{}, nil).Once()

			strategyRegistry.RegisterStrategy(mockStrat1)
			strategyRegistry.RegisterStrategy(mockStrat2)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Wait for MinExecutionInterval to pass (100ms default)
			time.Sleep(110 * time.Millisecond)

			// Trigger execution - need to reach threshold (default 5)
			for i := 0; i < 5; i++ {
				orchestrator.NotifyDataUpdate()
			}

			// Allow time for async execution
			time.Sleep(50 * time.Millisecond)
		})

		It("should handle strategy returning error", func() {
			// Setup mock strategy that returns error
			mockStrat := mockStrategy.NewStrategy(GinkgoT())
			mockStrat.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).Return(nil, errors.New("strategy error")).Once()

			strategyRegistry.RegisterStrategy(mockStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Wait for MinExecutionInterval to pass (100ms default)
			time.Sleep(110 * time.Millisecond)

			// Trigger execution - need to reach threshold (default 5)
			for i := 0; i < 5; i++ {
				orchestrator.NotifyDataUpdate()
			}

			// Should handle error gracefully
			time.Sleep(50 * time.Millisecond)
		})

		It("should handle no strategies registered", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			orchestrator.NotifyDataUpdate()

			// Should not panic
			time.Sleep(100 * time.Millisecond)
		})
	})

	Describe("Concurrent Execution Prevention", func() {
		It("should prevent concurrent execution of the same strategy", func() {
			// Setup slow strategy
			mockStrat := mockStrategy.NewStrategy(GinkgoT())
			mockStrat.EXPECT().GetName().Return(strategy.StrategyName("SlowStrategy")).Maybe()
			mockStrat.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).RunAndReturn(func(ctx strategy.StrategyContext) ([]*strategy.Signal, error) {
				time.Sleep(200 * time.Millisecond)
				return []*strategy.Signal{}, nil
			}).Maybe()

			strategyRegistry.RegisterStrategy(mockStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Trigger multiple times rapidly to hit threshold multiple times
			// This will queue up multiple ticks while strategy is executing
			for i := 0; i < 15; i++ {
				go orchestrator.NotifyDataUpdate()
				time.Sleep(5 * time.Millisecond)
			}

			// Wait for executions to complete
			time.Sleep(800 * time.Millisecond)

			// Mutex prevents concurrent execution - no panics or deadlocks
		})
	})

	Describe("Panic Recovery", func() {
		It("should recover from strategy panic", func() {
			// Setup panicking strategy and normal strategy
			panicStrat := mockStrategy.NewStrategy(GinkgoT())
			normalStrat := mockStrategy.NewStrategy(GinkgoT())

			panicStrat.EXPECT().GetName().Return(strategy.StrategyName("PanicStrategy")).Maybe()
			panicStrat.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).RunAndReturn(func(ctx strategy.StrategyContext) ([]*strategy.Signal, error) {
				panic("intentional panic for testing")
			}).Once()

			normalStrat.EXPECT().GetName().Return(strategy.StrategyName("NormalStrategy")).Maybe()
			normalStrat.EXPECT().GetSignals(mock.AnythingOfType("strategy.StrategyContext")).Return([]*strategy.Signal{}, nil).Once()

			strategyRegistry.RegisterStrategy(panicStrat)
			strategyRegistry.RegisterStrategy(normalStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Wait for MinExecutionInterval to pass (100ms default)
			time.Sleep(110 * time.Millisecond)

			// Trigger execution - need to reach threshold (default 5)
			for i := 0; i < 5; i++ {
				orchestrator.NotifyDataUpdate()
			}

			// Panic should be recovered
			time.Sleep(50 * time.Millisecond)
		})
	})

	Describe("GetStrategies", func() {
		It("should return all registered strategies", func() {
			mockStrat1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrat2 := mockStrategy.NewStrategy(GinkgoT())

			mockStrat1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat2.EXPECT().GetName().Return(strategy.StrategyName("Strategy2")).Maybe()

			strategyRegistry.RegisterStrategy(mockStrat1)
			strategyRegistry.RegisterStrategy(mockStrat2)

			strategies := orchestrator.GetStrategies()

			Expect(strategies).To(HaveLen(2))
		})

		It("should return empty list when no strategies registered", func() {
			strategies := orchestrator.GetStrategies()

			Expect(strategies).To(BeEmpty())
		})
	})
})
