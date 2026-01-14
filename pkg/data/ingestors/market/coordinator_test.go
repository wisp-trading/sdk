package market_test

import (
	"context"
	"time"

	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("MarketDataCoordinator", func() {
	var (
		app         *fxtest.App
		coordinator ingestors.MarketDataCoordinator
		ctx         context.Context
		cancel      context.CancelFunc
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&coordinator),
			fx.NopLogger,
		)
		Expect(app.Start(context.Background())).To(Succeed())
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		if coordinator != nil && coordinator.IsRunning() {
			_ = coordinator.StopDataCollection()
		}
		if cancel != nil {
			cancel()
		}
		if app != nil {
			Expect(app.Stop(context.Background())).To(Succeed())
		}
	})

	Describe("StartDataCollection", func() {

		Context("when starting data collection", func() {
			It("should start successfully and report running status", func() {
				err := coordinator.StartDataCollection(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(coordinator.IsRunning()).To(BeTrue())
			})
		})

		Context("when already running", func() {
			It("should be idempotent", func() {
				err := coordinator.StartDataCollection(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Try starting again - should not error
				err = coordinator.StartDataCollection(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(coordinator.IsRunning()).To(BeTrue())
			})
		})
	})

	Describe("StopDataCollection", func() {
		BeforeEach(func() {
			err := coordinator.StartDataCollection(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should stop all data collection", func() {
			err := coordinator.StopDataCollection()
			Expect(err).ToNot(HaveOccurred())
			Expect(coordinator.IsRunning()).To(BeFalse())
		})
	})

	Describe("GetStatus", func() {
		Context("when coordinator is not running", func() {
			It("should return status indicating not running", func() {
				status := coordinator.GetStatus()
				Expect(status).ToNot(BeNil())
				Expect(status["coordinator_running"]).To(BeFalse())
			})
		})

		Context("when coordinator is running", func() {
			BeforeEach(func() {
				err := coordinator.StartDataCollection(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return detailed status", func() {
				status := coordinator.GetStatus()
				Expect(status).ToNot(BeNil())
				Expect(status["coordinator_running"]).To(BeTrue())

				// Should have market types status
				marketTypes := status["market_types"]
				Expect(marketTypes).ToNot(BeNil())
			})
		})
	})

	Describe("ForceCollectNow", func() {
		BeforeEach(func() {
			err := coordinator.StartDataCollection(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should trigger immediate data collection without error", func() {
			// Should not panic or error
			coordinator.ForceCollectNow()
			// Give it a moment to process
			time.Sleep(100 * time.Millisecond)
		})
	})
})
