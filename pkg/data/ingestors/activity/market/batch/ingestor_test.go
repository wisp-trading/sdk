package batch_test

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

func setupMockConnector(t GinkgoTInterface, name connector.ExchangeName) *mockConnector.Connector {
	m := mockConnector.NewConnector(t)
	m.EXPECT().GetConnectorInfo().Return(&connector.Info{Name: name}).Maybe()
	m.EXPECT().SupportsPerpetuals().Return(true).Maybe()
	m.EXPECT().SupportsSpot().Return(false).Maybe()
	return m
}

var _ = Describe("BatchIngestor", func() {
	var (
		app               *fxtest.App
		batchIngestor     ingestors.BatchIngestor
		store             marketTypes.MarketData
		connectorRegistry registryTypes.ConnectorRegistry
		assetRegistry     registryTypes.AssetRegistry
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(&batchIngestor, &store, &connectorRegistry, &assetRegistry),
		)
		Expect(app.Start(context.Background())).To(Succeed())
	})

	AfterEach(func() {
		Expect(app.Stop(context.Background())).To(Succeed())
	})

	Describe("CollectNow", func() {
		Context("when fetching market data", func() {
			It("should collect orderbooks and klines for all assets", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockConnector(GinkgoT(), exchangeName)

				btc := portfolio.NewAsset("BTC")
				eth := portfolio.NewAsset("ETH")

				// Setup orderbook expectations
				btcOrderbook := &connector.OrderBook{
					Asset: btc,
					Bids:  []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:  []connector.PriceLevel{{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)}},
				}
				ethOrderbook := &connector.OrderBook{
					Asset: eth,
					Bids:  []connector.PriceLevel{{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(1.0)}},
					Asks:  []connector.PriceLevel{{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(1.0)}},
				}

				m.EXPECT().FetchOrderBook(btc, connector.TypePerpetual, 20).Return(btcOrderbook, nil).Maybe()
				m.EXPECT().FetchOrderBook(eth, connector.TypePerpetual, 20).Return(ethOrderbook, nil).Maybe()

				// Setup kline expectations for all intervals
				intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
				for _, interval := range intervals {
					m.EXPECT().FetchKlines("BTC", interval, mock.Anything).Return([]connector.Kline{
						{Symbol: "BTC", Interval: interval, Open: 50000, High: 50100, Low: 49900, Close: 50050, Volume: 100, OpenTime: time.Now().Add(-time.Hour), CloseTime: time.Now()},
					}, nil).Maybe()
					m.EXPECT().FetchKlines("ETH", interval, mock.Anything).Return([]connector.Kline{
						{Symbol: "ETH", Interval: interval, Open: 3000, High: 3010, Low: 2990, Close: 3005, Volume: 50, OpenTime: time.Now().Add(-time.Hour), CloseTime: time.Now()},
					}, nil).Maybe()
				}

				// Register connector and assets
				connectorRegistry.RegisterConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())
				assetRegistry.RegisterAsset(btc, connector.TypePerpetual)
				assetRegistry.RegisterAsset(eth, connector.TypePerpetual)

				// Act
				batchIngestor.CollectNow()
				time.Sleep(200 * time.Millisecond)

				// Assert - data should be stored
				btcOB := store.GetOrderBook(btc, exchangeName, connector.TypePerpetual)
				Expect(btcOB).ToNot(BeNil(), "BTC orderbook should be stored")

				ethOB := store.GetOrderBook(eth, exchangeName, connector.TypePerpetual)
				Expect(ethOB).ToNot(BeNil(), "ETH orderbook should be stored")

				btcKlines := store.GetKlines(btc, exchangeName, "1m", 10)
				Expect(btcKlines).ToNot(BeEmpty(), "BTC klines should be stored")
			})
		})

		Context("when no assets are required", func() {
			It("should skip collection without errors", func() {
				exchangeName := connector.ExchangeName("test-exchange")
				m := setupMockConnector(GinkgoT(), exchangeName)

				connectorRegistry.RegisterConnector(exchangeName, m)
				Expect(connectorRegistry.MarkConnectorReady(exchangeName)).To(Succeed())

				// Act - no assets registered, should complete without calling connector
				batchIngestor.CollectNow()
				time.Sleep(50 * time.Millisecond)
			})
		})
	})

	Describe("Start and Stop", func() {
		It("should start and stop correctly", func() {
			err := batchIngestor.Start(100 * time.Millisecond)
			Expect(err).ToNot(HaveOccurred())
			Expect(batchIngestor.IsActive()).To(BeTrue())

			err = batchIngestor.Stop()
			Expect(err).ToNot(HaveOccurred())
			Expect(batchIngestor.IsActive()).To(BeFalse())
		})
	})
})
