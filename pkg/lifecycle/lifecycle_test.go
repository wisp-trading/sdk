package lifecycle_test

import (
	"context"
	"time"

	mockSpotConnector "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/connector/spot"
	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	lifecycleTypes "github.com/backtesting-org/kronos-sdk/pkg/types/lifecycle"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	registryTypes "github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func setupMockSpotConnector(t GinkgoTInterface, name connector.ExchangeName) *mockSpotConnector.Connector {
	m := mockSpotConnector.NewConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()

	// Add expectations for batch ingestor calls (they will be called when data collection starts)
	m.EXPECT().FetchOrderBook(mock.Anything, mock.Anything).Return(&connector.OrderBook{}, nil).Maybe()
	m.EXPECT().FetchKlines(mock.Anything, mock.Anything, mock.Anything).Return([]connector.Kline{}, nil).Maybe()
	m.EXPECT().FetchPrice(mock.Anything).Return(&connector.Price{}, nil).Maybe()

	return m
}

var _ = Describe("LifecycleController", func() {
	var (
		app               *fxtest.App
		controller        lifecycleTypes.Controller
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.AssetRegistry
		ctx               context.Context
		cancel            context.CancelFunc
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&controller, &connectorRegistry, &assetRegistry),
			fx.NopLogger,
		)
		Expect(app.Start(context.Background())).To(Succeed())
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		if controller != nil && controller.IsReady() {
			_ = controller.Stop(ctx)
		}
		cancel()
		Expect(app.Stop(context.Background())).To(Succeed())
	})

	Describe("State Transitions", func() {
		Context("when starting the controller", func() {
			It("should transition from Created to Ready", func() {
				// Initial state
				Expect(controller.State()).To(Equal(lifecycleTypes.StateCreated))
				Expect(controller.IsReady()).To(BeFalse())

				// Register a connector so validation passes
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockSpotConnector(GinkgoT(), exchangeName)
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(portfolio.NewAsset("BTC"), connector.TypeSpot)

				// Start
				err := controller.Start(ctx, strategy.StrategyName("test-strategy"))
				Expect(err).ToNot(HaveOccurred())

				// Should be ready
				Expect(controller.State()).To(Equal(lifecycleTypes.StateReady))
				Expect(controller.IsReady()).To(BeTrue())
			})

			It("should fail if no connectors are registered", func() {
				// Try to start without any connectors
				err := controller.Start(ctx, strategy.StrategyName("test-strategy"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no connectors marked as ready"))
			})
		})

		Context("when stopping the controller", func() {
			BeforeEach(func() {
				// Register connector and start
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockSpotConnector(GinkgoT(), exchangeName)
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(portfolio.NewAsset("BTC"), connector.TypeSpot)

				err := controller.Start(ctx, strategy.StrategyName("test-strategy"))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should transition from Ready to Stopped", func() {
				// Verify started
				Expect(controller.State()).To(Equal(lifecycleTypes.StateReady))

				// Stop
				err := controller.Stop(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Should be stopped
				Expect(controller.State()).To(Equal(lifecycleTypes.StateStopped))
				Expect(controller.IsReady()).To(BeFalse())
			})

			It("should be idempotent", func() {
				// Stop once
				err := controller.Stop(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Stop again - should not error
				err = controller.Stop(ctx)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("WaitUntilReady", func() {
		It("should block until controller is ready", func() {
			// Register connector
			exchangeName := connector.ExchangeName("test-exchange")
			m := setupMockSpotConnector(GinkgoT(), exchangeName)
			connectorRegistry.RegisterSpotConnector(exchangeName, m)
			Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
			assetRegistry.RegisterAsset(portfolio.NewAsset("BTC"), connector.TypeSpot)

			// Start in background
			go func() {
				time.Sleep(100 * time.Millisecond)
				_ = controller.Start(ctx, strategy.StrategyName("test-strategy"))
			}()

			// Wait for ready with timeout
			waitCtx, waitCancel := context.WithTimeout(ctx, 1*time.Second)
			defer waitCancel()

			err := controller.WaitUntilReady(waitCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(controller.IsReady()).To(BeTrue())
		})

		It("should timeout if controller never becomes ready", func() {
			// Don't start the controller
			waitCtx, waitCancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer waitCancel()

			err := controller.WaitUntilReady(waitCtx)
			Expect(err).To(Equal(context.DeadlineExceeded))
		})
	})

	Describe("Start Validation", func() {
		Context("when trying to start multiple times", func() {
			It("should reject second start attempt", func() {
				// Register connector and start
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockSpotConnector(GinkgoT(), exchangeName)
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(portfolio.NewAsset("BTC"), connector.TypeSpot)

				// First start
				err := controller.Start(ctx, strategy.StrategyName("test-strategy"))
				Expect(err).ToNot(HaveOccurred())

				// Second start should fail
				err = controller.Start(ctx, strategy.StrategyName("test-strategy"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot start"))
			})
		})
	})
})
