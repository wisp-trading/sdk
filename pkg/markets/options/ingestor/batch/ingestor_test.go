package batch_test

import (
	"time"

	"github.com/stretchr/testify/mock"
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	optionsDomain "github.com/wisp-trading/sdk/pkg/markets/options"
	optionsBatch "github.com/wisp-trading/sdk/pkg/markets/options/ingestor/batch"
	optionsStore "github.com/wisp-trading/sdk/pkg/markets/options/store"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	mockConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/options"
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

func setupMockOptionsConnector(t GinkgoTInterface, name connector.ExchangeName) *mockConnector.MockConnector {
	m := mockConnector.NewMockConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Options BatchIngestor", func() {
	var (
		store             optionsTypes.OptionsStore
		connectorRegistry registryTypes.ConnectorRegistry
		watchlist         optionsTypes.OptionsWatchlist
		logger            logging.ApplicationLogger
		timeProviderInst  temporal.TimeProvider
		factory           batchTypes.BatchIngestorFactory
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

		factory = optionsBatch.NewFactory(
			connectorRegistry,
			watchlist,
			store,
			timeProviderInst,
			logger,
		)
	})

	Describe("Factory", func() {
		Context("when there are ready options connectors", func() {
			It("should create one ingestor per connector", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockOptionsConnector(GinkgoT(), exchangeName)
				m.EXPECT().GetExpirationData(mock.Anything, mock.Anything).Return(make(map[float64]map[string]optionsconnector.OptionData), nil).Maybe()

				connectorRegistry.RegisterOptions(exchangeName, m)
				connectorRegistry.MarkReady(exchangeName)

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))
				Expect(ingestors[0].GetMarketType()).To(Equal(connector.MarketTypeOptions))
			})
		})

		Context("when no connectors are ready", func() {
			It("should return empty ingestor list", func() {
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(BeEmpty())
			})
		})
	})

	Describe("Batch Ingestor Collection", func() {
		var (
			ingestor batchTypes.BatchIngestor
			m        *mockConnector.MockConnector
		)

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			m = setupMockOptionsConnector(GinkgoT(), exchangeName)

			connectorRegistry.RegisterOptions(exchangeName, m)
			connectorRegistry.MarkReady(exchangeName)

			ingestors := factory.CreateIngestors()
			Expect(ingestors).To(HaveLen(1))
			ingestor = ingestors[0]
		})

		Context("when collecting watched expirations", func() {
			It("should fetch and store option data", func() {
				watchlist.RequireExpiration(connector.ExchangeName("test-exchange"), btcPair, expiration)

				expirationData := map[float64]map[string]optionsconnector.OptionData{
					50000: {
						"CALL": {
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
						},
					},
				}
				m.EXPECT().GetExpirationData(btcPair, expiration).Return(expirationData, nil).Maybe()

				ingestor.CollectNow()

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
			})
		})
	})

	Describe("Batch Ingestor Lifecycle", func() {
		var (
			ingestor batchTypes.BatchIngestor
			m        *mockConnector.MockConnector
		)

		BeforeEach(func() {
			exchangeName := connector.ExchangeName("test-exchange")
			m = setupMockOptionsConnector(GinkgoT(), exchangeName)
			m.EXPECT().GetExpirationData(mock.Anything, mock.Anything).Return(make(map[float64]map[string]optionsconnector.OptionData), nil).Maybe()

			connectorRegistry.RegisterOptions(exchangeName, m)
			connectorRegistry.MarkReady(exchangeName)

			ingestors := factory.CreateIngestors()
			Expect(ingestors).To(HaveLen(1))
			ingestor = ingestors[0]
		})

		It("should start and stop correctly", func() {
			Expect(ingestor.IsActive()).To(BeFalse())

			err := ingestor.Start(100 * time.Millisecond)
			Expect(err).NotTo(HaveOccurred())
			Expect(ingestor.IsActive()).To(BeTrue())

			time.Sleep(50 * time.Millisecond)

			err = ingestor.Stop()
			Expect(err).NotTo(HaveOccurred())
			Expect(ingestor.IsActive()).To(BeFalse())
		})

		It("should have correct market type", func() {
			Expect(ingestor.GetMarketType()).To(Equal(connector.MarketTypeOptions))
		})
	})
})
