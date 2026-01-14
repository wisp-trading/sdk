package realtime_test

import (
	"context"
	"time"

	mockConnector "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	registryTypes "github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func setupMockWSConnector(t GinkgoTInterface, name connector.ExchangeName) *mockConnector.WebSocketConnector {
	m := mockConnector.NewWebSocketConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{Name: name}).Maybe()
	m.EXPECT().SupportsPerpetuals().Return(true).Maybe()
	m.EXPECT().SupportsSpot().Return(false).Maybe()
	return m
}

var _ = Describe("RealtimeIngestor", func() {
	var (
		app               *fxtest.App
		realtimeIngestor  ingestors.RealtimeIngestor
		store             marketTypes.MarketData
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.AssetRegistry
		ctx               context.Context
		cancel            context.CancelFunc
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&realtimeIngestor, &store, &connectorRegistry, &assetRegistry),
		)
		Expect(app.Start(context.Background())).To(Succeed())
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
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC-PERP": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

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
				}
				orderbookChan <- orderbook
				time.Sleep(200 * time.Millisecond)

				// Assert - data should be stored
				storedOB := store.GetOrderBook(btc, exchangeName, connector.TypePerpetual)
				Expect(storedOB).ToNot(BeNil(), "Orderbook should be stored")
				Expect(storedOB.Bids).To(HaveLen(1))
				Expect(storedOB.Asks).To(HaveLen(1))

				close(orderbookChan)
			})
		})

		Context("when receiving kline updates", func() {
			It("should process and store kline data", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				klineChan := make(chan connector.Kline, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

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

				close(klineChan)
			})
		})

		Context("when channel closes unexpectedly", func() {
			It("should handle closure gracefully without panic", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				klineChan := make(chan connector.Kline, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

				// Start ingestor
				Expect(realtimeIngestor.Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				// Close channel - should not panic
				close(klineChan)
				time.Sleep(100 * time.Millisecond)

				// Test passes if no panic
			})
		})
	})
})
