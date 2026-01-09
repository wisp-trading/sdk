package batch_test

import (
	stdtime "time"

	mockConnector "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	mockIngestors "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	mockMarket "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	mockHealth "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/health"
	mockRegistry "github.com/backtesting-org/kronos-sdk/mocks/github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/runtime/time"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	ingestorTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	loggingTypes "github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("BatchIngestor", func() {
	var (
		ingestor             ingestorTypes.BatchIngestor
		mockStore            *mockMarket.MarketData
		mockExchangeRegistry *mockRegistry.ConnectorRegistry
		mockAssetRegistry    *mockRegistry.AssetRegistry
		logger               loggingTypes.ApplicationLogger
		mockHealthStore      *mockHealth.CoordinatorHealthStore
		mockNotifier         *mockIngestors.DataUpdateNotifier
		mockedConnector      *mockConnector.Connector
	)

	BeforeEach(func() {
		mockStore = mockMarket.NewMarketData(GinkgoT())
		mockExchangeRegistry = mockRegistry.NewConnectorRegistry(GinkgoT())
		mockAssetRegistry = mockRegistry.NewAssetRegistry(GinkgoT())
		logger = loggingTypes.NewNoOpLogger()
		mockHealthStore = mockHealth.NewCoordinatorHealthStore(GinkgoT())
		mockNotifier = mockIngestors.NewDataUpdateNotifier(GinkgoT())
		mockedConnector = mockConnector.NewConnector(GinkgoT())

		ingestor = batch.NewBatchIngestor(
			mockStore,
			mockExchangeRegistry,
			mockAssetRegistry,
			logger,
			time.NewTimeProvider(),
			mockHealthStore,
			mockNotifier,
		)
	})

	Describe("CollectNow", func() {
		Context("when fetching market data", func() {
			It("should collect orderbooks and klines for all assets", func() {
				btcAsset := portfolio.NewAsset("BTC")
				ethAsset := portfolio.NewAsset("ETH")

				mockAssetRegistry.EXPECT().GetRequiredAssets().Return([]portfolio.Asset{
					btcAsset,
					ethAsset,
				}).Times(2) // Once for initial collection, once for CollectNow

				mockedConnector.EXPECT().GetConnectorInfo().Return(&connector.Info{Name: "hyperliquid"}).Maybe()
				mockedConnector.EXPECT().SupportsPerpetuals().Return(true).Maybe()
				mockedConnector.EXPECT().SupportsSpot().Return(false).Maybe()

				mockExchangeRegistry.EXPECT().GetReadyConnectors().Return([]connector.Connector{mockedConnector}).Times(2) // Once for initial, once for CollectNow

				// Orderbook expectations (2x for initial + CollectNow)
				btcOrderbook := &connector.OrderBook{
					Asset: btcAsset,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				ethOrderbook := &connector.OrderBook{
					Asset: ethAsset,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(1.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				mockedConnector.EXPECT().FetchOrderBook(btcAsset, connector.TypePerpetual, 20).Return(btcOrderbook, nil).Times(2)
				mockedConnector.EXPECT().FetchOrderBook(ethAsset, connector.TypePerpetual, 20).Return(ethOrderbook, nil).Times(2)

				mockStore.EXPECT().UpdateOrderBook(btcAsset, connector.ExchangeName("hyperliquid"), connector.TypePerpetual, *btcOrderbook).Times(2)
				mockStore.EXPECT().UpdateOrderBook(ethAsset, connector.ExchangeName("hyperliquid"), connector.TypePerpetual, *ethOrderbook).Times(2)

				mockHealthStore.EXPECT().RecordDataReceived(connector.ExchangeName("hyperliquid"), mock.Anything, mock.Anything, mock.Anything).Times(4)

				// Kline expectations for multiple intervals (2x for initial + CollectNow)
				// Use the actual configured limits from the ingestor
				klineLimits := map[string]int{
					"1m":  500,
					"5m":  300,
					"15m": 200,
					"1h":  168,
					"4h":  180,
					"1d":  90,
				}

				intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
				for _, interval := range intervals {
					limit := klineLimits[interval]

					btcKlines := []connector.Kline{
						{
							Symbol:    "BTC",
							Interval:  interval,
							Open:      numerical.NewFromFloat(50000),
							High:      numerical.NewFromFloat(50100),
							Low:       numerical.NewFromFloat(49900),
							Close:     numerical.NewFromFloat(50050),
							Volume:    numerical.NewFromFloat(100),
							OpenTime:  stdtime.Now().Add(-stdtime.Hour),
							CloseTime: stdtime.Now(),
						},
					}

					ethKlines := []connector.Kline{
						{
							Symbol:    "ETH",
							Interval:  interval,
							Open:      numerical.NewFromFloat(3000),
							High:      numerical.NewFromFloat(3010),
							Low:       numerical.NewFromFloat(2990),
							Close:     numerical.NewFromFloat(3005),
							Volume:    numerical.NewFromFloat(50),
							OpenTime:  stdtime.Now().Add(-stdtime.Hour),
							CloseTime: stdtime.Now(),
						},
					}

					mockedConnector.EXPECT().FetchKlines("BTC", interval, limit).Return(btcKlines, nil).Times(2)
					mockedConnector.EXPECT().FetchKlines("ETH", interval, limit).Return(ethKlines, nil).Times(2)

					for _, kline := range btcKlines {
						mockStore.EXPECT().UpdateKline(btcAsset, connector.ExchangeName("hyperliquid"), kline).Times(2)
					}

					for _, kline := range ethKlines {
						mockStore.EXPECT().UpdateKline(ethAsset, connector.ExchangeName("hyperliquid"), kline).Times(2)
					}
				}

				mockNotifier.EXPECT().Notify().Times(2)

				// Start the ingestor
				err := ingestor.Start(30 * stdtime.Second)
				Expect(err).ToNot(HaveOccurred())

				// Wait for initial collection
				stdtime.Sleep(100 * stdtime.Millisecond)

				// Execute CollectNow
				ingestor.CollectNow()

				// Give goroutines time to complete
				stdtime.Sleep(500 * stdtime.Millisecond)

				// Cleanup
				err = ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())

				// Verify
				mockStore.AssertExpectations(GinkgoT())
				mockedConnector.AssertExpectations(GinkgoT())
				mockHealthStore.AssertExpectations(GinkgoT())
				mockNotifier.AssertExpectations(GinkgoT())
			})
		})

		Context("when no assets are required", func() {
			It("should skip collection", func() {
				mockAssetRegistry.EXPECT().GetRequiredAssets().Return([]portfolio.Asset{}).Times(2)

				err := ingestor.Start(30 * stdtime.Second)
				Expect(err).ToNot(HaveOccurred())

				stdtime.Sleep(100 * stdtime.Millisecond)

				ingestor.CollectNow()
				stdtime.Sleep(100 * stdtime.Millisecond)

				err = ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())

				mockAssetRegistry.AssertExpectations(GinkgoT())
			})
		})

		Context("when kline fetch fails", func() {
			It("should continue with other intervals and assets", func() {
				btcAsset := portfolio.NewAsset("BTC")

				mockAssetRegistry.EXPECT().GetRequiredAssets().Return([]portfolio.Asset{btcAsset}).Times(2)
				mockedConnector.EXPECT().GetConnectorInfo().Return(&connector.Info{Name: "hyperliquid"}).Maybe()
				mockedConnector.EXPECT().SupportsPerpetuals().Return(true).Maybe()
				mockedConnector.EXPECT().SupportsSpot().Return(false).Maybe()
				mockExchangeRegistry.EXPECT().GetReadyConnectors().Return([]connector.Connector{mockedConnector}).Times(2)

				// Orderbook succeeds (2x)
				btcOrderbook := &connector.OrderBook{
					Asset: btcAsset,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				mockedConnector.EXPECT().FetchOrderBook(btcAsset, connector.TypePerpetual, 20).Return(btcOrderbook, nil).Times(2)
				mockStore.EXPECT().UpdateOrderBook(btcAsset, connector.ExchangeName("hyperliquid"), connector.TypePerpetual, *btcOrderbook).Times(2)
				mockHealthStore.EXPECT().RecordDataReceived(connector.ExchangeName("hyperliquid"), mock.Anything, mock.Anything, mock.Anything).Times(2)

				// First interval fails (2x) - using configured limit for 1m
				mockedConnector.EXPECT().FetchKlines("BTC", "1m", 500).Return(nil, nil).Times(2)

				// Rest succeed with empty results (2x) - using configured limits
				klineLimits := map[string]int{
					"5m":  300,
					"15m": 200,
					"1h":  168,
					"4h":  180,
					"1d":  90,
				}

				intervals := []string{"5m", "15m", "1h", "4h", "1d"}
				for _, interval := range intervals {
					limit := klineLimits[interval]
					mockedConnector.EXPECT().FetchKlines("BTC", interval, limit).Return([]connector.Kline{}, nil).Times(2)
				}

				mockNotifier.EXPECT().Notify().Times(2)

				err := ingestor.Start(30 * stdtime.Second)
				Expect(err).ToNot(HaveOccurred())

				stdtime.Sleep(100 * stdtime.Millisecond)

				ingestor.CollectNow()
				stdtime.Sleep(500 * stdtime.Millisecond)

				err = ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())

				mockStore.AssertExpectations(GinkgoT())
			})
		})
	})

	Describe("Start and Stop", func() {
		Context("when starting batch collection", func() {
			It("should start ticker and collect periodically", func() {
				mockAssetRegistry.EXPECT().GetRequiredAssets().Return([]portfolio.Asset{}).Maybe()
				mockNotifier.EXPECT().Notify().Maybe()

				err := ingestor.Start(30 * stdtime.Second)
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeTrue())

				// Stop
				err = ingestor.Stop()
				Expect(err).ToNot(HaveOccurred())
				Expect(ingestor.IsActive()).To(BeFalse())
			})
		})
	})
})
