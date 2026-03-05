package realtime_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	mockPerpConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/perp"
	perpDomain "github.com/wisp-trading/sdk/pkg/markets/perp"
	perpRealtime "github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/realtime"
	perpStore "github.com/wisp-trading/sdk/pkg/markets/perp/store"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/registry"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
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
		store             perpTypes.MarketStore
		connectorRegistry registryTypes.ConnectorRegistry
		watchlist         perpTypes.PerpWatchlist
		logger            logging.ApplicationLogger
		timeProviderInst  temporal.TimeProvider
		factory           realtimeTypes.RealtimeIngestorFactory
		ctx               context.Context
		cancel            context.CancelFunc
		btc               = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		timeProviderInst = timeProvider.NewTimeProvider()
		_ = timeProviderInst
		store = perpStore.NewStore(timeProvider.NewTimeProvider())
		connectorRegistry = registry.NewConnectorRegistry()
		watchlist = perpDomain.NewPerpWatchlist()

		factory = perpRealtime.NewFactory(
			connectorRegistry,
			watchlist,
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

				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				orderbookChannels := map[string]<-chan connector.OrderBook{"BTC": orderbookChan}
				klineChannels := map[string]<-chan connector.Kline{"BTC-1m": klineChan}

				m.EXPECT().GetOrderBookChannels().Return(orderbookChannels).Maybe()
				m.EXPECT().GetKlineChannels().Return(klineChannels).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				connectorRegistry.RegisterPerp(exchangeName, m)
				Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())
				watchlist.RequirePair(exchangeName, btc)

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				Expect(ingestors[0].Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				orderbookChan <- connector.OrderBook{
					Pair:      btc,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)}},
					Timestamp: time.Now(),
				}
				time.Sleep(200 * time.Millisecond)

				storedOB := store.GetOrderBook(btc, exchangeName)
				Expect(storedOB).ToNot(BeNil())
				Expect(storedOB.Bids[0].Price.InexactFloat64()).To(Equal(50000.0))
				Expect(storedOB.Asks[0].Price.InexactFloat64()).To(Equal(50100.0))

				close(orderbookChan)
				close(klineChan)
				close(fundingChan)
			})
		})

		Context("when receiving funding rate updates", func() {
			It("should process and store funding rate data (perp-specific)", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpWSConnector(GinkgoT(), exchangeName)

				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				m.EXPECT().GetOrderBookChannels().Return(map[string]<-chan connector.OrderBook{"BTC": orderbookChan}).Maybe()
				m.EXPECT().GetKlineChannels().Return(map[string]<-chan connector.Kline{"BTC-1m": klineChan}).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				connectorRegistry.RegisterPerp(exchangeName, m)
				Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())
				watchlist.RequirePair(exchangeName, btc)

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				Expect(ingestors[0].Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				now := time.Now()
				fundingChan <- perpConn.FundingRate{
					Pair:            btc,
					CurrentRate:     numerical.NewFromFloat(0.0001),
					NextFundingTime: now.Add(8 * time.Hour),
					MarkPrice:       numerical.NewFromFloat(50050),
					IndexPrice:      numerical.NewFromFloat(50045),
					Timestamp:       now,
				}
				time.Sleep(200 * time.Millisecond)

				storedFunding := store.GetFundingRate(btc, exchangeName)
				Expect(storedFunding).ToNot(BeNil())
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

				orderbookChan := make(chan connector.OrderBook, 10)
				klineChan := make(chan connector.Kline, 10)
				fundingChan := make(chan perpConn.FundingRate, 10)
				errorChan := make(chan error, 10)

				m.EXPECT().GetOrderBookChannels().Return(map[string]<-chan connector.OrderBook{"BTC": orderbookChan}).Maybe()
				m.EXPECT().GetKlineChannels().Return(map[string]<-chan connector.Kline{"BTC-1m": klineChan}).Maybe()
				m.EXPECT().FundingRateUpdates().Return((<-chan perpConn.FundingRate)(fundingChan)).Maybe()
				m.EXPECT().ErrorChannel().Return((<-chan error)(errorChan)).Maybe()
				m.EXPECT().StartWebSocket().Return(nil).Maybe()
				m.EXPECT().SubscribeOrderBook(mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeKlines(mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().SubscribeFundingRates(mock.Anything).Return(nil).Maybe()

				connectorRegistry.RegisterPerp(exchangeName, m)
				Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())
				watchlist.RequirePair(exchangeName, btc)

				ingestors := factory.CreateIngestors()
				Expect(ingestors[0].Start(ctx)).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				close(orderbookChan)
				close(klineChan)
				close(fundingChan)
				time.Sleep(100 * time.Millisecond)
				// Test passes if no panic
			})
		})
	})
})
