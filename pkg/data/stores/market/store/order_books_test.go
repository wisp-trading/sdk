package store_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store"
	timeProvider "github.com/wisp-trading/sdk/pkg/runtime/time"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("Market Data Store - OrderBooks", func() {
	var (
		marketStore marketTypes.MarketStore
		btc         portfolio.Pair
		eth         portfolio.Pair
		provider    temporal.TimeProvider
	)

	BeforeEach(func() {
		provider = timeProvider.NewTimeProvider()
		marketStore = store.NewStore(provider)
		btc = portfolio.NewPair(
			portfolio.NewAsset("BTC"),
			portfolio.NewAsset("USD"),
		)

		eth = portfolio.NewPair(
			portfolio.NewAsset("ETH"),
			portfolio.NewAsset("USD"),
		)
	})

	Describe("UpdateOrderBook", func() {
		Context("when adding a new orderbook", func() {
			It("should marketStore the orderbook correctly", func() {
				// Create an orderbook
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
						{Price: numerical.NewFromFloat(49990), Quantity: numerical.NewFromFloat(2.0)},
						{Price: numerical.NewFromFloat(49980), Quantity: numerical.NewFromFloat(3.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
						{Price: numerical.NewFromFloat(50020), Quantity: numerical.NewFromFloat(2.5)},
						{Price: numerical.NewFromFloat(50030), Quantity: numerical.NewFromFloat(3.5)},
					},
				}

				// Update the marketStore
				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook)

				// Retrieve and verify
				retrieved := marketStore.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(retrieved.Bids).To(HaveLen(3))
				Expect(retrieved.Asks).To(HaveLen(3))
				Expect(retrieved.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(retrieved.Asks[0].Price).To(Equal(numerical.NewFromFloat(50010)))
			})

			It("should handle multiple instrument types for the same asset and exchange", func() {
				now := time.Now()

				orderBookPerp := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				orderBookSpot := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50005), Quantity: numerical.NewFromFloat(2.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50015), Quantity: numerical.NewFromFloat(1.5)},
					},
				}

				// Note: In the new architecture, spot and perp connectors would be registered separately
				// For this test, we'll use different exchange names to simulate separate connectors
				marketStore.UpdateOrderBook(btc, "hyperliquid-perp", orderBookPerp)
				marketStore.UpdateOrderBook(btc, "hyperliquid-spot", orderBookSpot)

				perpBook := marketStore.GetOrderBook(btc, "hyperliquid-perp")
				spotBook := marketStore.GetOrderBook(btc, "hyperliquid-spot")

				Expect(perpBook).NotTo(BeNil())
				Expect(spotBook).NotTo(BeNil())
				Expect(perpBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(spotBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50005)))
			})

			It("should handle multiple exchanges for the same asset", func() {
				now := time.Now()

				orderBookHyper := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				orderBookBybit := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(2.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50110), Quantity: numerical.NewFromFloat(1.5)},
					},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBookHyper)
				marketStore.UpdateOrderBook(btc, "bybit", orderBookBybit)

				hyperBook := marketStore.GetOrderBook(btc, "hyperliquid")
				bybitBook := marketStore.GetOrderBook(btc, "bybit")

				Expect(hyperBook).NotTo(BeNil())
				Expect(bybitBook).NotTo(BeNil())
				Expect(hyperBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(bybitBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50100)))
			})

			It("should handle multiple assets", func() {
				now := time.Now()

				orderBookBTC := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				orderBookETH := connector.OrderBook{
					Pair:      eth,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(5.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(4.0)},
					},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBookBTC)
				marketStore.UpdateOrderBook(eth, "hyperliquid", orderBookETH)

				btcBook := marketStore.GetOrderBook(btc, "hyperliquid")
				ethBook := marketStore.GetOrderBook(eth, "hyperliquid")

				Expect(btcBook).NotTo(BeNil())
				Expect(ethBook).NotTo(BeNil())
				Expect(btcBook.Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(ethBook.Pair.Symbol()).To(Equal(eth.Symbol()))
				Expect(btcBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(ethBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(3000)))
			})
		})

		Context("when updating an existing orderbook", func() {
			It("should replace the orderbook with new data", func() {
				now := time.Now()

				// First orderbook
				orderBook1 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
						{Price: numerical.NewFromFloat(49990), Quantity: numerical.NewFromFloat(2.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
						{Price: numerical.NewFromFloat(50020), Quantity: numerical.NewFromFloat(2.5)},
					},
				}

				// Updated orderbook with different prices/quantities
				orderBook2 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now.Add(time.Second),
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50005), Quantity: numerical.NewFromFloat(2.0)},
						{Price: numerical.NewFromFloat(49995), Quantity: numerical.NewFromFloat(3.0)},
						{Price: numerical.NewFromFloat(49985), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50015), Quantity: numerical.NewFromFloat(1.5)},
					},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook1)
				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook2)

				retrieved := marketStore.GetOrderBook(btc, "hyperliquid")

				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Bids).To(HaveLen(3), "Should have updated bid levels")
				Expect(retrieved.Asks).To(HaveLen(1), "Should have updated ask levels")
				Expect(retrieved.Bids[0].Price).To(Equal(numerical.NewFromFloat(50005)), "Best bid should be updated")
				Expect(retrieved.Bids[0].Quantity).To(Equal(numerical.NewFromFloat(2.0)), "Best bid quantity should be updated")
				Expect(retrieved.Asks[0].Price).To(Equal(numerical.NewFromFloat(50015)), "Best ask should be updated")
			})
		})

		Context("when dealing with empty orderbooks", func() {
			It("should handle orderbooks with empty bids and asks", func() {
				now := time.Now()

				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{},
					Asks:      []connector.PriceLevel{},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook)

				retrieved := marketStore.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Bids).To(HaveLen(0))
				Expect(retrieved.Asks).To(HaveLen(0))
			})
		})

		Context("when dealing with concurrent updates", func() {
			It("should handle concurrent updates without data loss", func() {
				now := time.Now()
				done := make(chan bool)

				// Spawn 10 goroutines updating different assets/exchanges
				for i := 0; i < 10; i++ {
					go func(idx int) {
						defer GinkgoRecover()
						orderBook := connector.OrderBook{
							Pair:      btc,
							Timestamp: now,
							Bids: []connector.PriceLevel{
								{Price: numerical.NewFromFloat(50000 + float64(idx*10)), Quantity: numerical.NewFromFloat(1.0)},
							},
							Asks: []connector.PriceLevel{
								{Price: numerical.NewFromFloat(50100 + float64(idx*10)), Quantity: numerical.NewFromFloat(1.0)},
							},
						}
						marketStore.UpdateOrderBook(btc, connector.ExchangeName("exchange"+string(rune(idx))), orderBook)
						done <- true
					}(i)
				}

				// Wait for all goroutines to complete
				for i := 0; i < 10; i++ {
					<-done
				}

				// Verify we can retrieve orderbooks for all exchanges
				books := marketStore.GetOrderBooks(btc)
				Expect(books).To(HaveLen(10))
			})
		})
	})

	Describe("GetOrderBooks", func() {
		Context("when orderbooks exist for an asset", func() {
			It("should return all orderbooks for that asset", func() {
				now := time.Now()

				orderBook1 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)}},
				}

				orderBook2 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50100), Quantity: numerical.NewFromFloat(2.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50110), Quantity: numerical.NewFromFloat(1.5)}},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook1)
				marketStore.UpdateOrderBook(btc, "bybit", orderBook2)

				books := marketStore.GetOrderBooks(btc)
				Expect(books).To(HaveLen(2))
				Expect(books["hyperliquid"]).NotTo(BeNil())
				Expect(books["bybit"]).NotTo(BeNil())
			})
		})

		Context("when no orderbooks exist for an asset", func() {
			It("should return an empty map", func() {
				books := marketStore.GetOrderBooks(btc)
				Expect(books).To(HaveLen(0))
			})
		})
	})

	Describe("OrderBook", func() {
		Context("when orderbook exists", func() {
			It("should return the correct orderbook", func() {
				now := time.Now()

				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)}},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook)

				retrieved := marketStore.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Pair.Symbol()).To(Equal(btc.Symbol()))
			})
		})

		Context("when orderbook does not exist", func() {
			It("should return nil for non-existent asset", func() {
				retrieved := marketStore.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).To(BeNil())
			})

			It("should return nil for non-existent exchange", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)}},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook)

				retrieved := marketStore.GetOrderBook(btc, "bybit")
				Expect(retrieved).To(BeNil())
			})
		})
	})

	Describe("GetAllAssetsWithOrderBooks", func() {
		Context("when orderbooks exist", func() {
			It("should return all assets with orderbooks", func() {
				now := time.Now()

				orderBookBTC := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)}},
				}

				orderBookETH := connector.OrderBook{
					Pair:      eth,
					Timestamp: now,
					Bids:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3000), Quantity: numerical.NewFromFloat(5.0)}},
					Asks:      []connector.PriceLevel{{Price: numerical.NewFromFloat(3010), Quantity: numerical.NewFromFloat(4.0)}},
				}

				marketStore.UpdateOrderBook(btc, "hyperliquid", orderBookBTC)
				marketStore.UpdateOrderBook(eth, "hyperliquid", orderBookETH)

				assets := marketStore.GetAllAssetsWithOrderBooks()
				Expect(assets).To(HaveLen(2))
				Expect(assets).To(ContainElement(btc))
				Expect(assets).To(ContainElement(eth))
			})
		})

		Context("when no orderbooks exist", func() {
			It("should return an empty slice", func() {
				assets := marketStore.GetAllAssetsWithOrderBooks()
				Expect(assets).To(HaveLen(0))
			})
		})
	})

	Describe("Thread safety and data integrity", func() {
		It("should maintain data integrity under high concurrent load", func() {
			now := time.Now()
			iterations := 100
			done := make(chan bool, iterations*2)

			// Concurrent writers
			for i := 0; i < iterations; i++ {
				go func(idx int) {
					defer GinkgoRecover()
					orderBook := connector.OrderBook{
						Pair:      btc,
						Timestamp: now,
						Bids: []connector.PriceLevel{
							{Price: numerical.NewFromFloat(50000 + float64(idx)), Quantity: numerical.NewFromFloat(float64(idx))},
						},
						Asks: []connector.PriceLevel{
							{Price: numerical.NewFromFloat(50100 + float64(idx)), Quantity: numerical.NewFromFloat(float64(idx))},
						},
					}
					marketStore.UpdateOrderBook(btc, "hyperliquid", orderBook)
					done <- true
				}(i)
			}

			// Concurrent readers
			for i := 0; i < iterations; i++ {
				go func() {
					defer GinkgoRecover()
					_ = marketStore.GetOrderBook(btc, "hyperliquid")
					_ = marketStore.GetOrderBooks(btc)
					_ = marketStore.GetAllAssetsWithOrderBooks()
					done <- true
				}()
			}

			// Wait for all operations to complete
			for i := 0; i < iterations*2; i++ {
				<-done
			}

			// Verify final state is valid
			retrieved := marketStore.GetOrderBook(btc, "hyperliquid")
			Expect(retrieved).NotTo(BeNil())
			Expect(retrieved.Pair.Symbol()).To(Equal(btc.Symbol()))
		})
	})
})
