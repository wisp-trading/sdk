package realtime_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"

	realtime "github.com/wisp-trading/sdk/pkg/data/ingestors/market/prediction/real_time"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/data"
	predictionStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	mockPredictionConn "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/prediction"
)

var _ = Describe("Prediction Watchlist Integration", func() {
	var (
		app       *fxtest.App
		ctx       context.Context
		cancel    context.CancelFunc
		watchlist data.PredictionWatchlist
		store     predictionStore.MarketStore
		logger    logging.ApplicationLogger
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&watchlist, &store, &logger),
			fx.NopLogger,
		)

		Expect(app.Start(context.Background())).To(Succeed())
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
		if app != nil {
			Expect(app.Stop(context.Background())).To(Succeed())
		}
	})

	Describe("Factory and Ingestor Watchlist Subscription", func() {
		var (
			mockWSConn    *mockPredictionConn.WebSocketConnector
			exchangeName  connector.ExchangeName
			testMarket1   prediction.Market
			testMarket2   prediction.Market
			orderBookChan chan prediction.OrderBook
			extension     realtimeTypes.PredictionExtension
		)

		BeforeEach(func() {
			exchangeName = "test-exchange"
			testMarket1 = prediction.Market{
				MarketID:    "market-1",
				Slug:        "rain-tomorrow",
				Exchange:    exchangeName,
				OutcomeType: prediction.OutcomeTypeBinary,
				Active:      true,
				Outcomes: []prediction.Outcome{
					{OutcomeID: "yes"},
					{OutcomeID: "no"},
				},
			}
			testMarket2 = prediction.Market{
				MarketID:    "market-2",
				Slug:        "price-up",
				Exchange:    exchangeName,
				OutcomeType: prediction.OutcomeTypeBinary,
				Active:      true,
				Outcomes: []prediction.Outcome{
					{OutcomeID: "up"},
					{OutcomeID: "down"},
				},
			}

			// Create mock WebSocket connector
			mockWSConn = mockPredictionConn.NewWebSocketConnector(GinkgoT())
			orderBookChan = make(chan prediction.OrderBook, 100)

			// Setup connector expectations
			mockWSConn.EXPECT().GetConnectorInfo().Return(&connector.Info{
				Name:             exchangeName,
				WebSocketEnabled: true,
			}).Maybe()
			mockWSConn.EXPECT().StartWebSocket().Return(nil).Maybe()
			mockWSConn.EXPECT().StopWebSocket().Return(nil).Maybe()
			mockWSConn.EXPECT().ErrorChannel().Return(make(<-chan error)).Maybe()
			mockWSConn.EXPECT().OrderBookUpdates().Return((<-chan prediction.OrderBook)(orderBookChan)).Maybe()
			mockWSConn.EXPECT().SubscribeOrderBook(testMarket1).Return(nil).Maybe()
			mockWSConn.EXPECT().SubscribeOrderBook(testMarket2).Return(nil).Maybe()
			mockWSConn.EXPECT().UnsubscribeMarket(testMarket1).Return(nil).Maybe()
			mockWSConn.EXPECT().UnsubscribeMarket(testMarket2).Return(nil).Maybe()

			// Create orderbook extension
			extension = realtime.NewPredictionOrderBookExtension(
				store.(predictionStore.OrderBookStoreExtension),
				logger,
			)
		})

		AfterEach(func() {
			if orderBookChan != nil {
				close(orderBookChan)
			}
		})
		Context("when ingestor is created and started", func() {
			It("should subscribe to existing markets in watchlist", func() {
				// Arrange - add market to watchlist before starting ingestor
				watchlist.RequireMarket(exchangeName, testMarket1)

				// Act - create ingestor directly
				ingestor := realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())
				Expect(ingestor.GetMarketType()).To(Equal(connector.MarketTypePrediction))

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeTrue())

				// Give it time to process initial subscriptions
				time.Sleep(100 * time.Millisecond)

				// Assert - verify subscription was called
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)

				// Cleanup
				_ = ingestor.Stop()
			})

			It("should start successfully even with no markets in watchlist", func() {
				// Act - create ingestor directly with empty watchlist
				ingestor := realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeTrue())

				time.Sleep(50 * time.Millisecond)

				// Assert - no subscriptions should have been called
				mockWSConn.AssertNotCalled(GinkgoT(), "SubscribeOrderBook")

				// Cleanup
				_ = ingestor.Stop()
			})
		})

		Context("when markets are dynamically added to watchlist", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				// Create ingestor directly
				ingestor = realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeTrue())
			})

			AfterEach(func() {
				if ingestor != nil && ingestor.IsActive() {
					_ = ingestor.Stop()
				}
			})

			It("should subscribe to newly added markets", func() {
				// Act - add market to watchlist after ingestor is running
				watchlist.RequireMarket(exchangeName, testMarket1)

				// Give it time to process the watchlist event
				time.Sleep(150 * time.Millisecond)

				// Assert - subscription should have been called
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)
			})

			It("should subscribe to multiple markets added sequentially", func() {
				// Act - add multiple markets
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				watchlist.RequireMarket(exchangeName, testMarket2)
				time.Sleep(100 * time.Millisecond)

				// Assert - both subscriptions should have been called
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket2)
			})

			It("should handle adding the same market twice (idempotent)", func() {
				// Act - add same market twice
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				// Reset mock call count
				initialCalls := len(mockWSConn.Calls)

				// Add same market again - watchlist is idempotent
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				// Assert - watchlist should not emit duplicate event
				Expect(mockWSConn.Calls).To(HaveLen(initialCalls),
					"Should not subscribe again for duplicate market")
			})
		})

		Context("when markets are removed from watchlist", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				// Add markets first
				watchlist.RequireMarket(exchangeName, testMarket1)
				watchlist.RequireMarket(exchangeName, testMarket2)

				// Create ingestor directly
				ingestor = realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				time.Sleep(100 * time.Millisecond)
			})

			AfterEach(func() {
				if ingestor != nil && ingestor.IsActive() {
					_ = ingestor.Stop()
				}
			})

			It("should handle removing markets gracefully", func() {
				// Act - remove market
				watchlist.ReleaseMarket(exchangeName, testMarket1.MarketID)
				time.Sleep(150 * time.Millisecond)

				// Assert - should still be running
				Expect(ingestor.IsActive()).To(BeTrue())
			})

			It("should handle removing non-existent market gracefully", func() {
				// Act - remove market that wasn't added
				nonExistentMarketID := prediction.MarketID("non-existent")

				// Should not panic
				watchlist.ReleaseMarket(exchangeName, nonExistentMarketID)
				time.Sleep(50 * time.Millisecond)

				// Assert - should still be running
				Expect(ingestor.IsActive()).To(BeTrue())
			})
		})

		Context("when processing order book updates", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				// Add market and start ingestor
				watchlist.RequireMarket(exchangeName, testMarket1)

				// Create ingestor directly
				ingestor = realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				time.Sleep(100 * time.Millisecond)
			})

			AfterEach(func() {
				if ingestor != nil && ingestor.IsActive() {
					_ = ingestor.Stop()
				}
			})

			It("should store order book updates in the store", func() {
				// Arrange - create order book update
				orderBookUpdate := prediction.OrderBook{
					MarketID:  testMarket1.MarketID,
					OutcomeID: "yes",
					OrderBook: connector.OrderBook{
						Bids: []connector.PriceLevel{
							{Price: numerical.NewFromFloat(0.55), Quantity: numerical.NewFromFloat(100)},
							{Price: numerical.NewFromFloat(0.54), Quantity: numerical.NewFromFloat(200)},
						},
						Asks: []connector.PriceLevel{
							{Price: numerical.NewFromFloat(0.56), Quantity: numerical.NewFromFloat(150)},
							{Price: numerical.NewFromFloat(0.57), Quantity: numerical.NewFromFloat(250)},
						},
					},
				}

				// Act - send order book update
				orderBookChan <- orderBookUpdate
				time.Sleep(100 * time.Millisecond)

				// Assert - verify it's stored
				storedOB := store.GetOrderBook(
					exchangeName,
					testMarket1.MarketID,
					"yes",
				)
				Expect(storedOB).ToNot(BeNil())
				Expect(storedOB.Bids).To(HaveLen(2))
				Expect(storedOB.Asks).To(HaveLen(2))
				Expect(storedOB.Bids[0].Price.Float64()).To(Equal(0.55))
			})
		})

		Context("when stopping the ingestor", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				watchlist.RequireMarket(exchangeName, testMarket1)

				// Create ingestor directly
				ingestor = realtime.NewPredictionRealtimeIngestor(
					mockWSConn,
					exchangeName,
					connector.MarketTypePrediction,
					watchlist,
					logger,
					extension,
				)
				Expect(ingestor).ToNot(BeNil())

				err := ingestor.Start(ctx)
				Expect(err).ToNot(HaveOccurred())
				time.Sleep(100 * time.Millisecond)
			})

			It("should stop cleanly and unsubscribe from watchlist", func() {
				// Act - stop ingestor
				err := ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeFalse())

				// Assert - should have called StopWebSocket
				mockWSConn.AssertCalled(GinkgoT(), "StopWebSocket")
			})

			It("should not process new watchlist events after stopping", func() {
				// Stop ingestor
				err := ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())

				// Reset call tracking
				mockWSConn.Calls = nil

				// Add new market
				watchlist.RequireMarket(exchangeName, testMarket2)
				time.Sleep(100 * time.Millisecond)

				// Should not subscribe since ingestor is stopped
				mockWSConn.AssertNotCalled(GinkgoT(), "SubscribeOrderBook", testMarket2)
			})
		})
	})
})
