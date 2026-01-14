package realtime_test

import (
	"context"
	"time"

	mockPerpConnector "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/perp/realtime"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/registry"
	timeProvider "github.com/backtesting-org/kronos-sdk/pkg/runtime/time"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	perpConn "github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	realtimeTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/realtime"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	registryTypes "github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func setupMockPerpWSConnector(t GinkgoTInterface, name connector.ExchangeName) *mockPerpConnector.WebSocketConnector {
	m := mockPerpConnector.NewWebSocketConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Perp RealtimeIngestor", func() {
	var (
		store             perp.MarketStore
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.AssetRegistry
		logger            logging.ApplicationLogger
		timeProviderInst  temporal.TimeProvider
		factory           realtimeTypes.RealtimeIngestorFactory
		ctx               context.Context
		cancel            context.CancelFunc
	)

	BeforeEach(func() {
		// Create real instances
		logger = logging.NewNoOpLogger()
		timeProviderInst = timeProvider.NewTimeProvider()
		store = perpStore.NewStore(timeProviderInst)
		connectorRegistry = registry.NewConnectorRegistry()
		assetRegistry = registry.NewAssetRegistry()

		// Create factory
		factory = realtime.NewFactory(
			connectorRegistry,
			assetRegistry,
			store,
			logger,
		)

		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		time.Sleep(50 * time.Millisecond)
	})

	Describe("WebSocket data ingestion", func() {
		Context("when receiving orderbook updates", func() {
			It("should process and store orderbook data", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterPerpConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

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
				close(fundingChan)
			})
		})

		Context("when receiving kline updates", func() {
			It("should process and store kline data", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterPerpConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

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
				close(fundingChan)
			})
		})

		Context("when receiving funding rate updates", func() {
			It("should process and store funding rate data (perp-specific)", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterPerpConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

				// Create ingestors from factory
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				realtimeIngestor := ingestors[0]

				// Start ingestor
				Expect(realtimeIngestor.Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				// Send funding rate update
				now := time.Now()
				fundingRate := perpConn.FundingRate{
					Asset:           btc,
					CurrentRate:     numerical.NewFromFloat(0.0001),
					NextFundingTime: now.Add(8 * time.Hour),
					MarkPrice:       numerical.NewFromFloat(50050),
					IndexPrice:      numerical.NewFromFloat(50045),
					Timestamp:       now,
				}
				fundingChan <- fundingRate
				time.Sleep(200 * time.Millisecond)

				// Assert - data should be stored
				storedFunding := store.GetFundingRate(btc, exchangeName)
				Expect(storedFunding).ToNot(BeNil(), "Funding rate should be stored")
				Expect(storedFunding.CurrentRate.InexactFloat64()).To(Equal(0.0001))
				Expect(storedFunding.MarkPrice.InexactFloat64()).To(Equal(50050.0))

				close(orderbookChan)
				close(klineChan)
				close(fundingChan)
			})
		})

		Context("when channel closes unexpectedly", func() {
			It("should handle closure gracefully without panic", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpWSConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")

				// Setup channels
				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{
					"BTC": orderbookChan,
				}
				klineChannels := map[string]<-chan connector.Kline{
					"BTC-1m": klineChan,
				}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				// Register connector and assets
				connectorRegistry.RegisterPerpConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)

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
				close(fundingChan)
				time.Sleep(100 * time.Millisecond)

				// Test passes if no panic
			})
		})
	})
})
