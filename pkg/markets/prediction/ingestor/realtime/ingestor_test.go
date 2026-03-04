package realtime_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/ingestor/realtime"
	prediction "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	mockPredictionConn "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
)

var _ = Describe("Prediction Watchlist Integration", func() {
	var (
		app       *fxtest.App
		ctx       context.Context
		cancel    context.CancelFunc
		watchlist types.PredictionWatchlist
		store     types.MarketStore
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
			extension     types.PredictionExtension
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
			mockWSConn.EXPECT().GetOrderBookUpdates().Return((<-chan prediction.OrderBook)(orderBookChan)).Maybe()
			mockWSConn.EXPECT().SubscribeOrderBook(testMarket1).Return(nil).Maybe()
			mockWSConn.EXPECT().SubscribeOrderBook(testMarket2).Return(nil).Maybe()
			mockWSConn.EXPECT().UnsubscribeMarket(testMarket1).Return(nil).Maybe()
			mockWSConn.EXPECT().UnsubscribeMarket(testMarket2).Return(nil).Maybe()

			// Create orderbook extension
			extension = realtime.NewPredictionOrderBookExtension(
				store.(types.OrderBookStoreExtension),
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
				watchlist.RequireMarket(exchangeName, testMarket1)

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

				time.Sleep(100 * time.Millisecond)

				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)

				_ = ingestor.Stop()
			})

			It("should start successfully even with no markets in watchlist", func() {
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

				mockWSConn.AssertNotCalled(GinkgoT(), "SubscribeOrderBook")

				_ = ingestor.Stop()
			})
		})

		Context("when markets are dynamically added to watchlist", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
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
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(150 * time.Millisecond)
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)
			})

			It("should subscribe to multiple markets added sequentially", func() {
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				watchlist.RequireMarket(exchangeName, testMarket2)
				time.Sleep(100 * time.Millisecond)

				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket1)
				mockWSConn.AssertCalled(GinkgoT(), "SubscribeOrderBook", testMarket2)
			})

			It("should handle adding the same market twice (idempotent)", func() {
				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				initialCalls := len(mockWSConn.Calls)

				watchlist.RequireMarket(exchangeName, testMarket1)
				time.Sleep(100 * time.Millisecond)

				Expect(mockWSConn.Calls).To(HaveLen(initialCalls),
					"Should not subscribe again for duplicate market")
			})
		})

		Context("when markets are removed from watchlist", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				watchlist.RequireMarket(exchangeName, testMarket1)
				watchlist.RequireMarket(exchangeName, testMarket2)

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
				watchlist.ReleaseMarket(exchangeName, testMarket1.MarketID)
				time.Sleep(150 * time.Millisecond)
				Expect(ingestor.IsActive()).To(BeTrue())
			})

			It("should handle removing non-existent market gracefully", func() {
				nonExistentMarketID := prediction.MarketID("non-existent")
				watchlist.ReleaseMarket(exchangeName, nonExistentMarketID)
				time.Sleep(50 * time.Millisecond)
				Expect(ingestor.IsActive()).To(BeTrue())
			})
		})

		Context("when processing order book updates", func() {
			var ingestor realtimeTypes.RealtimeIngestor

			BeforeEach(func() {
				watchlist.RequireMarket(exchangeName, testMarket1)

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

				orderBookChan <- orderBookUpdate
				time.Sleep(100 * time.Millisecond)

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
				err := ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeFalse())
				mockWSConn.AssertCalled(GinkgoT(), "StopWebSocket")
			})

			It("should not process new watchlist events after stopping", func() {
				err := ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())

				mockWSConn.Calls = nil

				watchlist.RequireMarket(exchangeName, testMarket2)
				time.Sleep(100 * time.Millisecond)

				mockWSConn.AssertNotCalled(GinkgoT(), "SubscribeOrderBook", testMarket2)
			})
		})
	})
})
