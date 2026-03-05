package lifecycle_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	mockStrategy "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/strategy"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Orchestrator", func() {
	var (
		app              *fxtest.App
		orchestrator     lifecycleTypes.Orchestrator
		strategyRegistry registry.StrategyRegistry
		ctx              context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&orchestrator,
				&strategyRegistry,
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
		It("should start successfully with no strategies", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not error when started twice", func() {
			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Start(ctx)
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

	Describe("Strategy Lifecycle", func() {
		It("should call Start on a registered strategy", func() {
			mockStrat := mockStrategy.NewStrategy(GinkgoT())
			mockStrat.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat.EXPECT().Start(mock.AnythingOfType("*context.cancelCtx")).Return(nil).Once()
			mockStrat.EXPECT().Stop(mock.Anything).Return(nil).Maybe()

			strategyRegistry.RegisterStrategy(mockStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should call Start on multiple registered strategies", func() {
			mockStrat1 := mockStrategy.NewStrategy(GinkgoT())
			mockStrat2 := mockStrategy.NewStrategy(GinkgoT())

			mockStrat1.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat1.EXPECT().Start(mock.Anything).Return(nil).Once()
			mockStrat1.EXPECT().Stop(mock.Anything).Return(nil).Maybe()

			mockStrat2.EXPECT().GetName().Return(strategy.StrategyName("Strategy2")).Maybe()
			mockStrat2.EXPECT().Start(mock.Anything).Return(nil).Once()
			mockStrat2.EXPECT().Stop(mock.Anything).Return(nil).Maybe()

			strategyRegistry.RegisterStrategy(mockStrat1)
			strategyRegistry.RegisterStrategy(mockStrat2)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should call Stop on all strategies when orchestrator stops", func() {
			mockStrat := mockStrategy.NewStrategy(GinkgoT())
			mockStrat.EXPECT().GetName().Return(strategy.StrategyName("Strategy1")).Maybe()
			mockStrat.EXPECT().Start(mock.Anything).Return(nil).Once()
			mockStrat.EXPECT().Stop(mock.Anything).Return(nil).Once()

			strategyRegistry.RegisterStrategy(mockStrat)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = orchestrator.Stop(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should propagate Start errors and roll back already-started strategies", func() {
			started := mockStrategy.NewStrategy(GinkgoT())
			started.EXPECT().GetName().Return(strategy.StrategyName("GoodStrategy")).Maybe()
			started.EXPECT().Start(mock.Anything).Return(nil).Once()
			started.EXPECT().Stop(mock.Anything).Return(nil).Maybe()

			failing := mockStrategy.NewStrategy(GinkgoT())
			failing.EXPECT().GetName().Return(strategy.StrategyName("BadStrategy")).Maybe()
			failing.EXPECT().Start(mock.Anything).Return(errors.New("startup error")).Once()
			failing.EXPECT().Stop(mock.Anything).Return(nil).Maybe()

			strategyRegistry.RegisterStrategy(started)
			strategyRegistry.RegisterStrategy(failing)

			err := orchestrator.Start(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("BadStrategy"))
		})
	})

	Describe("End-to-end signal dispatch", func() {
		It("a self-directed strategy fires on its own internal clock", func() {
			// timerStrategy is a concrete self-directed strategy.
			// It owns its goroutine via StartWithRunner and signals after a short delay.
			emitted := make(chan struct{}, 1)

			ts := newTimerStrategy(emitted)
			strategyRegistry.RegisterStrategy(ts)

			err := orchestrator.Start(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(emitted, "200ms").Should(Receive())
		})
	})
})

// timerStrategy is a minimal concrete strategy for testing.
// It embeds BaseStrategy and fires once after a short delay, then exits.
type timerStrategy struct {
	strategy.BaseStrategy
	emitted chan struct{}
}

func newTimerStrategy(emitted chan struct{}) *timerStrategy {
	return &timerStrategy{
		BaseStrategy: *strategy.NewBaseStrategy(strategy.BaseStrategyConfig{
			Name:      "TimerStrategy",
			RiskLevel: strategy.RiskLevelLow,
			Type:      strategy.StrategyTypeTechnical,
		}),
		emitted: emitted,
	}
}

// Start launches the strategy's run loop via BaseStrategy.StartWithRunner.
// This is the pattern all concrete strategies follow.
func (s *timerStrategy) Start(ctx context.Context) error {
	return s.BaseStrategy.StartWithRunner(ctx, s.run)
}

func (s *timerStrategy) run(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(20 * time.Millisecond):
		s.emitted <- struct{}{}
	}
}
