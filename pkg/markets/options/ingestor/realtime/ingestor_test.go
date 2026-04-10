package realtime_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	optionsDomain "github.com/wisp-trading/sdk/pkg/markets/options"
	optionsRealtime "github.com/wisp-trading/sdk/pkg/markets/options/ingestor/realtime"
	optionsStore "github.com/wisp-trading/sdk/pkg/markets/options/store"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	mockWSConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/registry"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func setupMockWSConnector(t GinkgoTInterface, name connector.ExchangeName) *mockWSConnector.WebSocketConnector {
	m := mockWSConnector.NewWebSocketConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Options RealtimeIngestor", func() {
	var (
		store             optionsTypes.OptionsStore
		connectorRegistry registryTypes.ConnectorRegistry
		watchlist         optionsTypes.OptionsWatchlist
		logger            logging.ApplicationLogger
		timeProviderInst  temporal.TimeProvider
		factory           optionsTypes.OptionsRealtimeIngestorFactory
		quote             = portfolio.NewAsset("USDT")
		btcPair           = portfolio.NewPair(portfolio.NewAsset("BTC"), quote)
		expiration        = time.Now().AddDate(0, 0, 30)
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		timeProviderInst = timeProvider.NewTimeProvider()
		store = optionsStore.NewStore(timeProviderInst)
		connectorRegistry = registry.NewConnectorRegistry()
		watchlist = optionsDomain.NewOptionsWatchlist()

		factory = optionsRealtime.NewFactory(
			connectorRegistry,
			watchlist,
			store,
			timeProviderInst,
			logger,
		)
	})

	Describe("Factory", func() {
		Context("when there are WebSocket-ready options connectors", func() {
			It("should create one ingestor per WebSocket-ready connector", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockWSConnector(GinkgoT(), exchangeName)
				m.EXPECT().GetTradeChannels().Return(map[string]<-chan connector.Trade{}).Maybe()
				m.EXPECT().GetOrderBookChannels().Return(map[string]<-chan connector.OrderBook{}).Maybe()
				m.EXPECT().SubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
				m.EXPECT().UnsubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

				connectorRegistry.RegisterOptions(exchangeName, m)
				connectorRegistry.MarkReady(exchangeName)

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))
				Expect(ingestors[0].GetMarketType()).To(Equal(connector.MarketTypeOptions))
			})
		})
	})

	Describe("Realtime Ingestor Subscriptions", func() {
		var (
			ingestor realtimeTypes.RealtimeIngestor
			m        *mockWSConnector.WebSocketConnector
			ctx      context.Context
			cancel   context.CancelFunc
		)

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			m = setupMockWSConnector(GinkgoT(), exchangeName)
			m.EXPECT().GetTradeChannels().Return(map[string]<-chan connector.Trade{}).Maybe()
			m.EXPECT().GetOrderBookChannels().Return(map[string]<-chan connector.OrderBook{}).Maybe()
			m.EXPECT().SubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			m.EXPECT().UnsubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

			connectorRegistry.RegisterOptions(exchangeName, m)
			connectorRegistry.MarkReady(exchangeName)

			ingestors := factory.CreateIngestors()
			Expect(ingestors).To(HaveLen(1))
			ingestor = ingestors[0]

			ctx, cancel = context.WithCancel(context.Background())
			DeferCleanup(cancel)
		})

		Context("when subscribing to watched expirations", func() {
			It("should subscribe to watched expirations on start", func() {
				m.EXPECT().GetOptionUpdateChannels().Return(map[string]<-chan optionsconnector.OptionUpdate{}).Maybe()
				watchlist.RequireExpiration(connector.ExchangeName("test-exchange"), btcPair, expiration)

				// Pre-populate watchlist strikes so the ingestor can resolve contracts
				watchlist.SetStrikes(connector.ExchangeName("test-exchange"), btcPair, expiration, []float64{50000})

				err := ingestor.Start(ctx)
				Expect(err).NotTo(HaveOccurred())

				m.AssertCalled(GinkgoT(), "SubscribeExpirationUpdates", btcPair, expiration, mock.Anything)

				err = ingestor.Stop()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when unsubscribing from expirations", func() {
			It("should unsubscribe from expirations on stop", func() {
				m.EXPECT().GetOptionUpdateChannels().Return(map[string]<-chan optionsconnector.OptionUpdate{}).Maybe()
				watchlist.RequireExpiration(connector.ExchangeName("test-exchange"), btcPair, expiration)

				// Pre-populate watchlist strikes so the ingestor can resolve contracts
				watchlist.SetStrikes(connector.ExchangeName("test-exchange"), btcPair, expiration, []float64{50000})

				err := ingestor.Start(ctx)
				Expect(err).NotTo(HaveOccurred())

				err = ingestor.Stop()
				Expect(err).NotTo(HaveOccurred())

				m.AssertCalled(GinkgoT(), "UnsubscribeExpirationUpdates", btcPair, expiration, mock.Anything)
			})
		})

		Context("when receiving real-time updates", func() {
			It("should store option updates from WebSocket", func() {
				optionChanImpl := make(chan optionsconnector.OptionUpdate, 10)
				optionChannels := map[string]<-chan optionsconnector.OptionUpdate{"options": optionChanImpl}
				m.EXPECT().GetOptionUpdateChannels().Return(optionChannels).Maybe()

				watchlist.RequireExpiration(connector.ExchangeName("test-exchange"), btcPair, expiration)

				err := ingestor.Start(ctx)
				Expect(err).NotTo(HaveOccurred())

				time.Sleep(50 * time.Millisecond)

				update := optionsconnector.OptionUpdate{
					Contract: optionsconnector.OptionContract{
						Pair:       btcPair,
						Strike:     50000,
						Expiration: expiration,
						OptionType: "CALL",
					},
					MarkPrice:       2000,
					UnderlyingPrice: 50000,
					IV:              0.25,
					Greeks: optionsconnector.Greeks{
						Delta: 0.6,
						Gamma: 0.01,
						Theta: -0.05,
						Vega:  5.0,
						Rho:   0.15,
					},
					Timestamp: time.Now(),
				}

				optionChanImpl <- update

				time.Sleep(100 * time.Millisecond)

				contract := optionsTypes.OptionContract{
					Pair:       btcPair,
					Strike:     50000,
					Expiration: expiration,
					OptionType: "CALL",
				}

				Expect(store.GetMarkPrice(contract)).To(Equal(2000.0))
				Expect(store.GetIV(contract)).To(Equal(0.25))

				greeks := store.GetGreeks(contract)
				Expect(greeks.Delta).To(Equal(0.6))
				Expect(greeks.Gamma).To(Equal(0.01))
				Expect(greeks.Theta).To(Equal(-0.05))

				close(optionChanImpl)

				err = ingestor.Stop()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Realtime Ingestor Lifecycle", func() {
		var (
			ingestor realtimeTypes.RealtimeIngestor
			m        *mockWSConnector.WebSocketConnector
			ctx      context.Context
			cancel   context.CancelFunc
		)

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			m = setupMockWSConnector(GinkgoT(), exchangeName)
			m.EXPECT().GetOptionUpdateChannels().Return(map[string]<-chan optionsconnector.OptionUpdate{}).Maybe()
			m.EXPECT().GetTradeChannels().Return(map[string]<-chan connector.Trade{}).Maybe()
			m.EXPECT().GetOrderBookChannels().Return(map[string]<-chan connector.OrderBook{}).Maybe()
			m.EXPECT().SubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			m.EXPECT().UnsubscribeExpirationUpdates(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

			connectorRegistry.RegisterOptions(exchangeName, m)
			connectorRegistry.MarkReady(exchangeName)

			ingestors := factory.CreateIngestors()
			Expect(ingestors).To(HaveLen(1))
			ingestor = ingestors[0]

			ctx, cancel = context.WithCancel(context.Background())
			DeferCleanup(cancel)
		})

		It("should start and stop correctly", func() {
			m.EXPECT().GetOptionUpdateChannels().Return(map[string]<-chan optionsconnector.OptionUpdate{}).Maybe()
			Expect(ingestor.IsActive()).To(BeFalse())

			err := ingestor.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(ingestor.IsActive()).To(BeTrue())

			err = ingestor.Stop()
			Expect(err).NotTo(HaveOccurred())
			Expect(ingestor.IsActive()).To(BeFalse())
		})

		It("should have correct market type", func() {
			Expect(ingestor.GetMarketType()).To(Equal(connector.MarketTypeOptions))
		})

		It("should cancel processing on context cancellation", func() {
			m.EXPECT().GetOptionUpdateChannels().Return(map[string]<-chan optionsconnector.OptionUpdate{}).Maybe()
			watchlist.RequireExpiration(connector.ExchangeName("test-exchange"), btcPair, expiration)
			watchlist.SetStrikes(connector.ExchangeName("test-exchange"), btcPair, expiration, []float64{50000})

			err := ingestor.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(ingestor.IsActive()).To(BeTrue())

			cancel()

			time.Sleep(50 * time.Millisecond)

			err = ingestor.Stop()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
