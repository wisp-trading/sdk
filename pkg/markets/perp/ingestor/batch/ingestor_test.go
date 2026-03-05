package batch_test

import (
	"time"

	mockPerpConnector "github.com/wisp-trading/sdk/mocks/github.com/wisp-trading/sdk/pkg/types/connector/perp"
	perpDomain "github.com/wisp-trading/sdk/pkg/markets/perp"
	perpBatch "github.com/wisp-trading/sdk/pkg/markets/perp/ingestor/batch"
	perpStore "github.com/wisp-trading/sdk/pkg/markets/perp/store"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/registry"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	batchTypes "github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func setupMockPerpConnector(t GinkgoTInterface, name connector.ExchangeName) *mockPerpConnector.Connector {
	m := mockPerpConnector.NewConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{
		Name: name,
	}).Maybe()
	return m
}

var _ = Describe("Perp BatchIngestor", func() {
	var (
		store             perpTypes.MarketStore
		connectorRegistry registryTypes.ConnectorRegistry
		watchlist         perpTypes.PerpWatchlist
		logger            logging.ApplicationLogger
		timeProviderInst  temporal.TimeProvider
		factory           batchTypes.BatchIngestorFactory
		quote             = portfolio.NewAsset("USD")
		btcPair           = portfolio.NewPair(portfolio.NewAsset("BTC"), quote)
		ethPair           = portfolio.NewPair(portfolio.NewAsset("ETH"), quote)
	)

	BeforeEach(func() {
		logger = logging.NewNoOpLogger()
		timeProviderInst = timeProvider.NewTimeProvider()
		store = perpStore.NewStore(timeProviderInst)
		connectorRegistry = registry.NewConnectorRegistry()
		watchlist = perpDomain.NewPerpWatchlist()

		factory = perpBatch.NewFactory(
			connectorRegistry,
			watchlist,
			store,
			timeProviderInst,
			logger,
		)
	})

	Describe("CollectNow", func() {
		Context("when fetching perp market data", func() {
			It("should collect orderbooks, klines, and funding rates for all assets", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpConnector(GinkgoT(), exchangeName)

				now := time.Now()

				btcOrderbook := &connector.OrderBook{
					Pair:      btcPair,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)}},
					Timestamp: now,
				}
				ethOrderbook := &connector.OrderBook{
					Pair:      ethPair,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(1.0)}},
					Timestamp: now,
				}

				m.EXPECT().FetchOrderBook(btcPair, 20).Return(btcOrderbook, nil).Maybe()
				m.EXPECT().FetchOrderBook(ethPair, 20).Return(ethOrderbook, nil).Maybe()

				btcPrice := &connector.Price{
					Symbol:    "BTC",
					Price:     numerical.NewFromFloat(50050),
					BidPrice:  numerical.NewFromFloat(50000),
					AskPrice:  numerical.NewFromFloat(50100),
					Source:    exchangeName,
					Timestamp: now,
				}
				ethPrice := &connector.Price{
					Symbol:    "ETH",
					Price:     numerical.NewFromFloat(3005),
					BidPrice:  numerical.NewFromFloat(3000),
					AskPrice:  numerical.NewFromFloat(3010),
					Source:    exchangeName,
					Timestamp: now,
				}
				m.EXPECT().FetchPrice(btcPair).Return(btcPrice, nil).Maybe()
				m.EXPECT().FetchPrice(ethPair).Return(ethPrice, nil).Maybe()

				intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
				for _, interval := range intervals {
					m.EXPECT().FetchKlines(btcPair, interval, mock.Anything).Return([]connector.Kline{
						{Symbol: "BTC", Interval: interval, Open: 50000, High: 50100, Low: 49900, Close: 50050, Volume: 100, OpenTime: now.Add(-time.Hour), CloseTime: now},
					}, nil).Maybe()
					m.EXPECT().FetchKlines(ethPair, interval, mock.Anything).Return([]connector.Kline{
						{Symbol: "ETH", Interval: interval, Open: 3000, High: 3010, Low: 2990, Close: 3005, Volume: 50, OpenTime: now.Add(-time.Hour), CloseTime: now},
					}, nil).Maybe()
				}

				btcFundingRate := perpConn.FundingRate{
					Pair: btcPair, CurrentRate: numerical.NewFromFloat(0.0001),
					NextFundingTime: now.Add(8 * time.Hour), MarkPrice: numerical.NewFromFloat(50050),
					IndexPrice: numerical.NewFromFloat(50045), Timestamp: now,
				}
				ethFundingRate := perpConn.FundingRate{
					Pair: ethPair, CurrentRate: numerical.NewFromFloat(0.00005),
					NextFundingTime: now.Add(8 * time.Hour), MarkPrice: numerical.NewFromFloat(3005),
					IndexPrice: numerical.NewFromFloat(3004), Timestamp: now,
				}

				m.EXPECT().FetchFundingRate(btcPair).Return(&btcFundingRate, nil).Maybe()
				m.EXPECT().FetchFundingRate(ethPair).Return(&ethFundingRate, nil).Maybe()

				allFundingRates := map[portfolio.Pair]perpConn.FundingRate{
					btcPair: btcFundingRate,
					ethPair: ethFundingRate,
				}
				m.EXPECT().FetchCurrentFundingRates().Return(allFundingRates, nil).Maybe()

				connectorRegistry.RegisterPerp(exchangeName, m)
				Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())
				watchlist.RequirePair(exchangeName, btcPair)
				watchlist.RequirePair(exchangeName, ethPair)

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1), "Should create one ingestor for the registered connector")

				batchIngestor := ingestors[0]
				batchIngestor.CollectNow()
				time.Sleep(300 * time.Millisecond)

				btcOB := store.GetOrderBook(btcPair, exchangeName)
				Expect(btcOB).ToNot(BeNil(), "BTC orderbook should be stored")
				Expect(btcOB.Bids).To(HaveLen(1))
				Expect(btcOB.Bids[0].Price.InexactFloat64()).To(Equal(50000.0))

				ethOB := store.GetOrderBook(ethPair, exchangeName)
				Expect(ethOB).ToNot(BeNil(), "ETH orderbook should be stored")
				Expect(ethOB.Asks).To(HaveLen(1))
				Expect(ethOB.Asks[0].Price.InexactFloat64()).To(Equal(3010.0))

				btcKlines := store.GetKlines(btcPair, exchangeName, "1m", 10)
				Expect(btcKlines).ToNot(BeEmpty(), "BTC klines should be stored")
				Expect(btcKlines[0].Symbol).To(Equal("BTC"))

				ethKlines := store.GetKlines(ethPair, exchangeName, "5m", 10)
				Expect(ethKlines).ToNot(BeEmpty(), "ETH klines should be stored")

				btcFunding := store.GetFundingRate(btcPair, exchangeName)
				Expect(btcFunding).ToNot(BeNil(), "BTC funding rate should be stored")
				Expect(btcFunding.CurrentRate.InexactFloat64()).To(Equal(0.0001))

				ethFunding := store.GetFundingRate(ethPair, exchangeName)
				Expect(ethFunding).ToNot(BeNil(), "ETH funding rate should be stored")
				Expect(ethFunding.CurrentRate.InexactFloat64()).To(Equal(0.00005))

				storedBtcPrice := store.GetPairPrice(btcPair, exchangeName)
				Expect(storedBtcPrice).ToNot(BeNil(), "BTC price should be stored")
				Expect(storedBtcPrice.Price.InexactFloat64()).To(BeNumerically("~", 50050.0, 1.0))
			})
		})

		Context("when no assets are required", func() {
			It("should skip collection without errors", func() {
				exchangeName := connector.ExchangeName("test-perp-exchange")
				m := setupMockPerpConnector(GinkgoT(), exchangeName)

				connectorRegistry.RegisterPerp(exchangeName, m)
				Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())

				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				batchIngestor := ingestors[0]
				batchIngestor.CollectNow()
				time.Sleep(50 * time.Millisecond)
			})
		})
	})

	Describe("Start and Stop", func() {
		It("should start and stop correctly", func() {
			exchangeName := connector.ExchangeName("test-perp-exchange")
			m := setupMockPerpConnector(GinkgoT(), exchangeName)

			connectorRegistry.RegisterPerp(exchangeName, m)
			Expect(connectorRegistry.MarkReady(exchangeName)).To(Succeed())

			ingestors := factory.CreateIngestors()
			Expect(ingestors).To(HaveLen(1))

			batchIngestor := ingestors[0]

			err := batchIngestor.Start(100 * time.Millisecond)
			Expect(err).ToNot(HaveOccurred())
			Expect(batchIngestor.IsActive()).To(BeTrue())

			err = batchIngestor.Stop()
			Expect(err).ToNot(HaveOccurred())
			Expect(batchIngestor.IsActive()).To(BeFalse())
		})
	})
})
