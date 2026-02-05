package realtime_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	mockSpotConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/spot"
	spotRealtime "github.com/wisp-trading/sdk/pkg/data/ingestors/market/spot/realtime"
	sdkTesting "github.com/wisp-trading/sdk/pkg/testing"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	spotTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/spot"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func setupMockSpotWSConnector(t GinkgoTInterface, name connector.ExchangeName) *mockSpotConnector.WebSocketConnector {
	m := mockSpotConnector.NewWebSocketConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Spot RealtimeIngestor", func() {
	var (
		app               *fxtest.App
		store             spotTypes.MarketStore
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.PairRegistry
		factory           realtime.RealtimeIngestorFactory
		logger            logging.ApplicationLogger
		ctx               context.Context
		cancel            context.CancelFunc
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				fx.Annotate(&store, fx.ParamTags(`name:"spot_market_store"`)),
				&connectorRegistry,
				&assetRegistry,
				&logger,
			),
			fx.NopLogger,
		)
		Expect(app.Start(context.Background())).To(Succeed())

		// Create factory manually since it's now in a group, not available by name
		factory = spotRealtime.NewFactory(connectorRegistry, assetRegistry, store, logger)

		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		time.Sleep(50 * time.Millisecond)
		Expect(app.Stop(context.Background())).To(Succeed())
	})

	Describe("WebSocket data ingestion", func() {
		Context("when receiving orderbook updates", func() {
			It("should process and store orderbook data", func() {
				exchangeName := connector.ExchangeName("test-spot-exchange")
				m := setupMockSpotWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterPair(btc, connector.TypeSpot)

				// Create ingestors from factory
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				realtimeIngestor := ingestors[0]

				// Start ingestor
				Expect(realtimeIngestor.Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				// Send orderbook update
				orderbook := connector.OrderBook{
					Asset: btc,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)},
					},
					Timestamp: time.Now(),
				}
				orderbookChan <- orderbook
				time.Sleep(200 * time.Millisecond)

				// Assert - data should be stored
				storedOB := store.GetOrderBook(btc, exchangeName)
				Expect(storedOB).ToNot(BeNil(), "Orderbook should be stored")
				Expect(storedOB.Bids).To(HaveLen(1))
				Expect(storedOB.Bids[0].Price.InexactFloat64()).To(Equal(50000.0))
				Expect(storedOB.Asks).To(HaveLen(1))
				Expect(storedOB.Asks[0].Price.InexactFloat64()).To(Equal(50100.0))

				close(orderbookChan)
				close(klineChan)
			})
		})

		Context("when receiving kline updates", func() {
			It("should process and store kline data", func() {
				exchangeName := connector.ExchangeName("test-spot-exchange")
				m := setupMockSpotWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterPair(btc, connector.TypeSpot)

				// Create ingestors from factory
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				realtimeIngestor := ingestors[0]

				// Start ingestor
				Expect(realtimeIngestor.Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				// Send kline update
				kline := connector.Kline{
					Symbol:    "BTC",
					Interval:  "1m",
					Open:      50000,
					High:      50100,
					Low:       49900,
					Close:     50050,
					Volume:    100,
					OpenTime:  time.Now().Add(-time.Minute),
					CloseTime: time.Now(),
				}
				klineChan <- kline
				time.Sleep(200 * time.Millisecond)

				// Assert - data should be stored
				storedKlines := store.GetKlines(btc, exchangeName, "1m", 10)
				Expect(storedKlines).ToNot(BeEmpty(), "Klines should be stored")
				Expect(storedKlines[0].Symbol).To(Equal("BTC"))
				Expect(storedKlines[0].Close).To(Equal(50050.0))

				close(orderbookChan)
				close(klineChan)
			})
		})

		Context("when channel closes unexpectedly", func() {
			It("should handle closure gracefully without panic", func() {
				exchangeName := connector.ExchangeName("test-spot-exchange")
				m := setupMockSpotWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterPair(btc, connector.TypeSpot)

				// Create ingestors from factory
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				realtimeIngestor := ingestors[0]

				// Start ingestor
				Expect(realtimeIngestor.Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				// Close channels - should not panic
				close(orderbookChan)
				close(klineChan)
				time.Sleep(100 * time.Millisecond)

				// Test passes if no panic
			})
		})
	})
})
