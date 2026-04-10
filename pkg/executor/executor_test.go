package executor_test

import (
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/wisp-trading/sdk/pkg/executor"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/strategy"

	mockOptions "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/options/types"
	mockPerp "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/perp/types"
	mockPred "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	mockSpot "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/spot/types"
	mockExecution "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/execution"
	mockLogging "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/logging"
	mockRegistry "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/registry"
	mockStrategy "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/strategy"
	mockTemporal "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/temporal"
)

var _ = Describe("Executor", func() {
	var (
		logger       *mockLogging.ApplicationLogger
		timeProvider *mockTemporal.TimeProvider
		hookRegistry *mockRegistry.Hooks
		execRecords  *mockExecution.ExecutionRecords
		spotExec     *mockSpot.SignalExecutor
		perpExec     *mockPerp.SignalExecutor
		predExec     *mockPred.SignalExecutor
		optionsExec  *mockOptions.SignalExecutor
		exec         execution.Executor
		signal       *mockSpot.SpotSignal
		now          time.Time
	)

	BeforeEach(func() {
		logger = mockLogging.NewApplicationLogger(GinkgoT())
		timeProvider = mockTemporal.NewTimeProvider(GinkgoT())
		hookRegistry = mockRegistry.NewHooks(GinkgoT())
		execRecords = mockExecution.NewExecutionRecords(GinkgoT())
		spotExec = mockSpot.NewSignalExecutor(GinkgoT())
		perpExec = mockPerp.NewSignalExecutor(GinkgoT())
		predExec = mockPred.NewSignalExecutor(GinkgoT())
		optionsExec = mockOptions.NewSignalExecutor(GinkgoT())

		now = time.Now()
		timeProvider.EXPECT().Now().Return(now).Maybe()
		logger.EXPECT().Info(mock.Anything, mock.Anything).Maybe()
		logger.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
		// ExecutionRecords.Add is called on every execution
		execRecords.EXPECT().Add(mock.Anything).Maybe()

		exec = executor.NewExecutor(
			logger, timeProvider, hookRegistry, execRecords,
			spotExec, perpExec, predExec, optionsExec,
		)

		signal = mockSpot.NewSpotSignal(GinkgoT())
		signal.EXPECT().GetID().Return(uuid.New()).Maybe()
		signal.EXPECT().GetStrategy().Return(strategy.StrategyName("test-strategy")).Maybe()
		signal.EXPECT().GetTimestamp().Return(now).Maybe()
	})

	Describe("ExecuteSignalWithResult", func() {
		Context("when execution succeeds with no hooks", func() {
			BeforeEach(func() {
				hookRegistry.EXPECT().GetHooks().Return([]execution.ExecutionHook{})
				spotExec.EXPECT().ExecuteSpotSignal(signal, mock.Anything, mock.Anything).
					Run(func(sig spotTypes.SpotSignal, ctx *execution.ExecutionContext, result *execution.ExecutionResult) {
						result.OrderIDs = []string{"order-1"}
					}).
					Return(nil)
			})

			It("returns Success=true with the order IDs", func() {
				result, err := exec.ExecuteSignalWithResult(signal)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Success).To(BeTrue())
				Expect(result.OrderIDs).To(ConsistOf("order-1"))
				Expect(result.Error).To(BeNil())
				Expect(result.HookError).To(BeNil())
			})
		})

		Context("when a BeforeExecute hook blocks execution", func() {
			var hookErr error

			BeforeEach(func() {
				hookErr = errors.New("risk limit exceeded")
				hook := mockExecution.NewExecutionHook(GinkgoT())
				hook.EXPECT().BeforeExecute(mock.Anything).Return(hookErr)
				hook.EXPECT().OnError(mock.Anything, hookErr).Return(nil)

				hookRegistry.EXPECT().GetHooks().Return([]execution.ExecutionHook{hook})
				logger.EXPECT().Warn(mock.Anything, mock.Anything).Maybe()
			})

			It("returns Success=false with the hook error and never calls the domain executor", func() {
				result, err := exec.ExecuteSignalWithResult(signal)

				Expect(err).To(MatchError(hookErr))
				Expect(result.Success).To(BeFalse())
				Expect(result.Error).To(MatchError(hookErr))
				// spotExec was never called — mock AssertExpectations enforces this
			})
		})

		Context("when the domain executor fails", func() {
			var execErr error

			BeforeEach(func() {
				execErr = errors.New("exchange unavailable")
				hook := mockExecution.NewExecutionHook(GinkgoT())
				hook.EXPECT().BeforeExecute(mock.Anything).Return(nil)
				hook.EXPECT().OnError(mock.Anything, execErr).Return(nil)

				hookRegistry.EXPECT().GetHooks().Return([]execution.ExecutionHook{hook})
				spotExec.EXPECT().ExecuteSpotSignal(signal, mock.Anything, mock.Anything).Return(execErr)
			})

			It("returns Success=false with the execution error", func() {
				result, err := exec.ExecuteSignalWithResult(signal)

				Expect(err).To(MatchError(execErr))
				Expect(result.Success).To(BeFalse())
				Expect(result.Error).To(MatchError(execErr))
				Expect(result.HookError).To(BeNil())
			})
		})

		Context("when an AfterExecute hook fails", func() {
			var hookErr error

			BeforeEach(func() {
				hookErr = errors.New("position tracking failed")

				hook := mockExecution.NewExecutionHook(GinkgoT())
				hook.EXPECT().BeforeExecute(mock.Anything).Return(nil)
				hook.EXPECT().AfterExecute(mock.Anything, mock.Anything).Return(hookErr)

				hookRegistry.EXPECT().GetHooks().Return([]execution.ExecutionHook{hook})
				spotExec.EXPECT().ExecuteSpotSignal(signal, mock.Anything, mock.Anything).Return(nil)
				logger.EXPECT().Error(mock.Anything, mock.Anything).Maybe()
			})

			It("sets HookError and marks Success=false — order was placed but post-execution state is inconsistent", func() {
				result, err := exec.ExecuteSignalWithResult(signal)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Success).To(BeFalse())
				Expect(result.HookError).To(MatchError(hookErr))
				Expect(result.Error).To(BeNil())
			})
		})
	})

	Describe("ExecuteSignal", func() {
		Context("when an AfterExecute hook fails", func() {
			BeforeEach(func() {
				hookErr := errors.New("audit hook failed")
				hook := mockExecution.NewExecutionHook(GinkgoT())
				hook.EXPECT().BeforeExecute(mock.Anything).Return(nil)
				hook.EXPECT().AfterExecute(mock.Anything, mock.Anything).Return(hookErr)

				hookRegistry.EXPECT().GetHooks().Return([]execution.ExecutionHook{hook})
				spotExec.EXPECT().ExecuteSpotSignal(signal, mock.Anything, mock.Anything).Return(nil)
				// ExecuteSignalWithResult logs AfterExecute failure; ExecuteSignal logs HookError.
				// Both use variadic args — the mock passes them as (format, argsSlice) = 2 args.
				logger.EXPECT().Error(mock.Anything, mock.Anything).Times(2)
			})

			It("returns nil (order placed) but logs the hook failure", func() {
				err := exec.ExecuteSignal(signal)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

var _ = Describe("ExecutorRouter", func() {
	var (
		mockExec *mockExecution.Executor
		logger   *mockLogging.ApplicationLogger
		router   execution.SignalRouter
		signal   *mockStrategy.Signal
	)

	BeforeEach(func() {
		mockExec = mockExecution.NewExecutor(GinkgoT())
		logger = mockLogging.NewApplicationLogger(GinkgoT())
		logger.EXPECT().Error(mock.Anything, mock.Anything).Maybe()

		router = executor.NewExecutorRouter(mockExec, logger)

		signal = mockStrategy.NewSignal(GinkgoT())
		signal.EXPECT().GetID().Return(uuid.New()).Maybe()
		signal.EXPECT().GetStrategy().Return(strategy.StrategyName("test-strategy")).Maybe()
	})

	Describe("RouteWithResult", func() {
		It("sends the ExecutionResult to the provided channel when execution succeeds", func() {
			expected := execution.ExecutionResult{
				Success:  true,
				OrderIDs: []string{"order-42"},
			}
			mockExec.EXPECT().ExecuteSignalWithResult(signal).Return(expected, nil)

			ch := make(chan execution.ExecutionResult, 1)
			router.RouteWithResult(signal, ch)

			var result execution.ExecutionResult
			Eventually(ch).Should(Receive(&result))
			Expect(result.Success).To(BeTrue())
			Expect(result.OrderIDs).To(ConsistOf("order-42"))
		})

		It("sends the result even when execution fails", func() {
			execErr := errors.New("order rejected")
			failed := execution.ExecutionResult{Success: false, Error: execErr}
			mockExec.EXPECT().ExecuteSignalWithResult(signal).Return(failed, execErr)

			ch := make(chan execution.ExecutionResult, 1)
			router.RouteWithResult(signal, ch)

			var result execution.ExecutionResult
			Eventually(ch).Should(Receive(&result))
			Expect(result.Success).To(BeFalse())
		})

		It("sends the result when an AfterExecute hook fails", func() {
			hookErr := errors.New("position tracking failed")
			withHookErr := execution.ExecutionResult{
				Success:   false,
				HookError: hookErr,
				OrderIDs:  []string{"order-99"},
			}
			mockExec.EXPECT().ExecuteSignalWithResult(signal).Return(withHookErr, nil)

			ch := make(chan execution.ExecutionResult, 1)
			router.RouteWithResult(signal, ch)

			var result execution.ExecutionResult
			Eventually(ch).Should(Receive(&result))
			Expect(result.Success).To(BeFalse())
			Expect(result.HookError).To(MatchError(hookErr))
			Expect(result.OrderIDs).To(ConsistOf("order-99"))
		})
	})

	Describe("Route", func() {
		It("dispatches fire-and-forget and does not block", func() {
			done := make(chan struct{})
			mockExec.EXPECT().ExecuteSignal(signal).
				Run(func(strategy.Signal) { close(done) }).
				Return(nil)

			router.Route(signal)

			Eventually(done).Should(BeClosed())
		})
	})
})
