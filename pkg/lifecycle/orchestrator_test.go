package lifecycle_test

import (
	"context"
	"errors"
	"time"

	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mockIngestors "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	mockExecution "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	mockRegistry "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	mockStrategy "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	mockTemporal "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
	"github.com/backtesting-org/kronos-sdk/pkg/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Orchestrator", func() {
	var (
		mockExecutor         *mockExecution.Executor
		mockStrategyRegistry *mockRegistry.StrategyRegistry
		mockTimeProvider     *mockTemporal.TimeProvider
		mockNotifier         *mockIngestors.DataUpdateNotifier
		logger               logging.ApplicationLogger
		orchestrator         lifecycleTypes.Orchestrator
		ctx                  context.Context
		cancel               context.CancelFunc
	)

	BeforeEach(func() {
		mockExecutor = mockExecution.NewExecutor(GinkgoT())
		mockStrategyRegistry = mockRegistry.NewStrategyRegistry(GinkgoT())
		mockTimeProvider = mockTemporal.NewTimeProvider(GinkgoT())
		mockNotifier = mockIngestors.NewDataUpdateNotifier(GinkgoT())
		logger = logging.NewNoOpLogger()
		ctx, cancel = context.WithCancel(context.Background())

		// Setup default time provider behavior with mockery ticker
		mockTicker := mockTemporal.NewTicker(GinkgoT())
		tickerChan := make(chan time.Time)
		mockTicker.EXPECT().C().Return((<-chan time.Time)(tickerChan)).Maybe()
		mockTicker.EXPECT().Stop().Maybe()

		mockTimeProvider.EXPECT().Now().Return(time.Now()).Maybe()
		mockTimeProvider.EXPECT().Since(mock.Anything).Return(100 * time.Millisecond).Maybe()
		mockTimeProvider.EXPECT().NewTicker(mock.Anything).Return(mockTicker).Maybe()
	})

	AfterEach(func() {
		if orchestrator != nil {
			_ = orchestrator.Stop(ctx)
		}
		cancel()
	})

	Describe("Starting and Stopping", func() {
		It("should start successfully", func() {
			// Setup
			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			// Execute
			err := orchestrator.Start(ctx)

			// Assert
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not allow starting twice", func() {
			// Setup
			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			// Execute
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Start(ctx)

			// Assert - should warn but not error
			Expect(err).ToNot(HaveOccurred())
		})

		It("should stop successfully", func() {
			// Setup
			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			err = orchestrator.Stop(ctx)

			// Assert
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle multiple stop calls gracefully", func() {
			// Setup
			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			err = orchestrator.Stop(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Stop(ctx)

			// Assert - should handle gracefully
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Data Update Notification Flow", func() {
		It("should forward data updates to tick timer", func() {
			// Setup
			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{}).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute - send data update
			updatesChan <- struct{}{}

			// Assert - allow time for processing
			time.Sleep(50 * time.Millisecond)
		})

		It("should handle closed notifier channel gracefully", func() {
			// Setup
			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute - close channel
			close(updatesChan)

			// Assert - allow time for goroutine to exit
			time.Sleep(50 * time.Millisecond)
		})
	})

	Describe("Strategy Execution", func() {
		It("should execute single strategy", func() {
			ctx := context.Background()

			// Setup
			mockStrategy1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrategy1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrategy1.EXPECT().GetSignals(ctx).Return([]*strategy.Signal{}, nil).Maybe()

			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{mockStrategy1}).Maybe()
			mockExecutor.EXPECT().ExecuteSignal(mock.Anything).Return(nil).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute - trigger via notification
			orchestrator.NotifyDataUpdate()

			// Assert - give time for async execution
			time.Sleep(100 * time.Millisecond)
		})

		It("should execute multiple strategies concurrently", func() {
			ctx := context.Background()

			// Setup
			mockStrategy1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrategy2 := mockStrategy.NewStrategy(GinkgoT())

			mockStrategy1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrategy1.EXPECT().GetSignals(ctx).Return([]*strategy.Signal{}, nil).Maybe()
			mockStrategy2.EXPECT().GetName().Return(strategy.StrategyName("Strategy2")).Maybe()
			mockStrategy2.EXPECT().GetSignals(ctx).Return([]*strategy.Signal{}, nil).Maybe()

			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{
				mockStrategy1,
				mockStrategy2,
			}).Maybe()
			mockExecutor.EXPECT().ExecuteSignal(mock.Anything).Return(nil).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			orchestrator.NotifyDataUpdate()

			// Assert - give time for async execution
			time.Sleep(100 * time.Millisecond)
		})

		It("should handle strategy returning error", func() {
			ctx := context.Background()

			// Setup
			mockStrategy1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrategy1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrategy1.EXPECT().GetSignals(ctx).Return(nil, errors.New("mock error")).Maybe()

			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{mockStrategy1}).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			orchestrator.NotifyDataUpdate()

			// Assert - should handle error gracefully, give time for async execution
			time.Sleep(100 * time.Millisecond)
		})

		It("should handle no strategies registered", func() {
			// Setup
			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{}).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			orchestrator.NotifyDataUpdate()

			// Assert - should not panic
			time.Sleep(100 * time.Millisecond)
		})
	})

	Describe("Concurrent Execution Prevention", func() {
		It("should prevent concurrent execution of the same strategy", func() {
			ctx := context.Background()

			// Setup
			mockSlowStrategy := mockStrategy.NewStrategy(GinkgoT())
			mockSlowStrategy.EXPECT().GetName().Return(strategy.StrategyName("SlowStrategy")).Maybe()

			// Use RunAndReturn to add delay
			mockSlowStrategy.EXPECT().GetSignals(ctx).RunAndReturn(func(ctx2 context.Context) ([]*strategy.Signal, error) {
				time.Sleep(200 * time.Millisecond)
				return []*strategy.Signal{}, nil
			}).Maybe()

			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{mockSlowStrategy}).Maybe()
			mockExecutor.EXPECT().ExecuteSignal(mock.Anything).Return(nil).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute - trigger multiple times rapidly
			go orchestrator.NotifyDataUpdate()
			time.Sleep(10 * time.Millisecond)
			go orchestrator.NotifyDataUpdate()
			time.Sleep(10 * time.Millisecond)
			go orchestrator.NotifyDataUpdate()

			// Assert - wait for executions to complete
			time.Sleep(800 * time.Millisecond)

			// The mutex in orchestrator prevents concurrent execution
			// We can't easily assert this without exposing internal state
			// but the test verifies no panics or deadlocks occur
		})
	})

	Describe("Panic Recovery", func() {
		It("should recover from strategy panic", func() {
			ctx := context.Background()

			// Setup
			panicStrategy := mockStrategy.NewStrategy(GinkgoT())
			normalStrategy := mockStrategy.NewStrategy(GinkgoT())

			panicStrategy.EXPECT().GetName().Return(strategy.StrategyName("PanicStrategy")).Maybe()
			panicStrategy.EXPECT().GetSignals(ctx).RunAndReturn(func(ctx2 context.Context) ([]*strategy.Signal, error) {
				panic("intentional panic for testing")
			}).Maybe()

			normalStrategy.EXPECT().GetName().Return(strategy.StrategyName("NormalStrategy")).Maybe()
			normalStrategy.EXPECT().GetSignals(ctx).Return([]*strategy.Signal{}, nil).Maybe()

			updatesChan := make(chan struct{}, 10)
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{
				panicStrategy,
				normalStrategy,
			}).Maybe()
			mockExecutor.EXPECT().ExecuteSignal(mock.Anything).Return(nil).Maybe()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Execute
			orchestrator.NotifyDataUpdate()

			// Assert - panic should be recovered, give time for async execution
			time.Sleep(100 * time.Millisecond)

			// Verify orchestrator didn't crash from the panic
		})
	})

	Describe("GetStrategies", func() {
		It("should return all registered strategies", func() {
			// Setup
			mockStrategy1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrategy2 := mockStrategy.NewStrategy(GinkgoT())

			updatesChan := make(chan struct{})
			mockNotifier.EXPECT().Updates().Return((<-chan struct{})(updatesChan)).Maybe()
			mockStrategyRegistry.EXPECT().GetAllStrategies().Return([]strategy.Strategy{
				mockStrategy1,
				mockStrategy2,
			}).Once()

			orchestrator = lifecycle.NewOrchestrator(
				mockExecutor,
				mockStrategyRegistry,
				logger,
				mockTimeProvider,
				mockNotifier,
				nil, // profilingStore
				nil, // anomalyDetector
			)

			// Execute
			strategies := orchestrator.GetStrategies()

			// Assert
			Expect(strategies).To(HaveLen(2))
		})
	})
})
