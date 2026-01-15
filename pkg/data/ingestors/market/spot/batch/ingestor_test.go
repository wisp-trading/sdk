package batch_test

import (
	"context"
	"time"

	mockSpotConnector "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/connector/spot"
	spotBatch "github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/spot/batch"
	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	spotTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/spot"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	registryTypes "github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
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
	return m
}

var _ = Describe("Spot BatchIngestor", func() {
	var (
		app               *fxtest.App
		store             spotTypes.MarketStore
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.AssetRegistry
		factory           *spotBatch.Factory
		timeProvider      temporal.TimeProvider
		logger            logging.ApplicationLogger
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				fx.Annotate(&store, fx.ParamTags(`name:"spot_market_store"`)),
				&connectorRegistry,
				&assetRegistry,
				&timeProvider,
				&logger,
			),
			fx.NopLogger,
		)
		Expect(app.Start(context.Background())).To(Succeed())

		factory = spotBatch.NewFactory(connectorRegistry, assetRegistry, store, timeProvider, logger)
	})

	AfterEach(func() {
		Expect(app.Stop(context.Background())).To(Succeed())
	})

	Describe("CollectNow", func() {
		Context("when fetching spot market data", func() {
			It("should collect orderbooks, klines, and prices for all assets", func() {
				exchangeName := connector.ExchangeName("test-spot-exchange")
				m := setupMockSpotConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")
				eth := portfolio.NewAsset("ETH")

				now := time.Now()

				// Setup orderbook expectations
				btcOrderbook := &connector.OrderBook{
					Asset:     btc,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)}},
					Timestamp: now,
				}
				ethOrderbook := &connector.OrderBook{
					Asset:     eth,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(1.0)}},
					Timestamp: now,
				}

				m.EXPECT().FetchOrderBook(btc, 20).Return(btcOrderbook, nil).Maybe()
				m.EXPECT().FetchOrderBook(eth, 20).Return(ethOrderbook, nil).Maybe()

				// Setup price expectations
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
				m.EXPECT().FetchPrice("BTC").Return(btcPrice, nil).Maybe()
				m.EXPECT().FetchPrice("ETH").Return(ethPrice, nil).Maybe()

				// Setup kline expectations for all intervals
				intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
				for _, interval := range intervals {
					m.EXPECT().FetchKlines("BTC", interval, mock.Anything).Return([]connector.Kline{
						{
							Symbol:    "BTC",
							Interval:  interval,
							Open:      50000,
							High:      50100,
							Low:       49900,
							Close:     50050,
							Volume:    100,
							OpenTime:  now.Add(-time.Hour),
							CloseTime: now,
						},
					}, nil).Maybe()
					m.EXPECT().FetchKlines("ETH", interval, mock.Anything).Return([]connector.Kline{
						{
							Symbol:    "ETH",
							Interval:  interval,
							Open:      3000,
							High:      3010,
							Low:       2990,
							Close:     3005,
							Volume:    50,
							OpenTime:  now.Add(-time.Hour),
							CloseTime: now,
						},
					}, nil).Maybe()
				}

				// Register connector and assets
				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypeSpot)
				assetRegistry.RegisterAsset(eth, connector.TypeSpot)

				// Create ingestors from factory
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1), "Should create one ingestor for the registered connector")

				batchIngestor := ingestors[0]

				// Act - trigger collection
				batchIngestor.CollectNow()
				time.Sleep(300 * time.Millisecond) // Give it time to complete

				// Assert - verify orderbooks are stored
				btcOB := store.GetOrderBook(btc, exchangeName)
				Expect(btcOB).ToNot(BeNil(), "BTC orderbook should be stored")
				Expect(btcOB.Bids).To(HaveLen(1))
				Expect(btcOB.Bids[0].Price.InexactFloat64()).To(Equal(50000.0))

				ethOB := store.GetOrderBook(eth, exchangeName)
				Expect(ethOB).ToNot(BeNil(), "ETH orderbook should be stored")
				Expect(ethOB.Asks).To(HaveLen(1))
				Expect(ethOB.Asks[0].Price.InexactFloat64()).To(Equal(3010.0))

				// Assert - verify klines are stored
				btcKlines := store.GetKlines(btc, exchangeName, "1m", 10)
				Expect(btcKlines).ToNot(BeEmpty(), "BTC klines should be stored")
				Expect(btcKlines[0].Symbol).To(Equal("BTC"))

				ethKlines := store.GetKlines(eth, exchangeName, "5m", 10)
				Expect(ethKlines).ToNot(BeEmpty(), "ETH klines should be stored")

				// Assert - verify prices are stored
				storedBtcPrice := store.GetAssetPrice(btc, exchangeName)
				Expect(storedBtcPrice).ToNot(BeNil(), "BTC price should be stored")
				Expect(storedBtcPrice.Price.InexactFloat64()).To(BeNumerically("~", 50050.0, 1.0))
			})
		})

		Context("when no assets are required", func() {
			It("should skip collection without errors", func() {
				exchangeName := connector.ExchangeName("test-spot-exchange")
				m := setupMockSpotConnector(GinkgoT(), exchangeName)

				connectorRegistry.RegisterSpotConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())

				// Create ingestors without registering any assets
				ingestors := factory.CreateIngestors()
				Expect(ingestors).To(HaveLen(1))

				batchIngestor := ingestors[0]

				// Act - no assets registered, should complete without calling connector
				batchIngestor.CollectNow()
				time.Sleep(50 * time.Millisecond)

				// Should not panic or error
			})
		})
	})

	Describe("Start and Stop", func() {
		It("should start and stop correctly", func() {
			exchangeName := connector.ExchangeName("test-spot-exchange")
			m := setupMockSpotConnector(GinkgoT(), exchangeName)

			connectorRegistry.RegisterSpotConnector(exchangeName, m)
			Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())

			// Create ingestor
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
