package extensions_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("OrderBookExtension - Unit Tests", func() {
	type priceUpdate struct {
		pair     portfolio.Pair
		exchange connector.ExchangeName
		price    connector.Price
	}

	var (
		extension           marketTypes.OrderBookStoreExtension
		btc                 portfolio.Pair
		eth                 portfolio.Pair
		priceUpdates        []priceUpdate
		metadataUpdates     []marketTypes.UpdateKey
		priceUpdateMutex    sync.Mutex
		metadataUpdateMutex sync.Mutex
	)

	// Helper to capture price updates
	capturePriceUpdate := func(pair portfolio.Pair, exchange connector.ExchangeName, price connector.Price) {
		priceUpdateMutex.Lock()
		defer priceUpdateMutex.Unlock()
		priceUpdates = append(priceUpdates, priceUpdate{
			pair:     pair,
			exchange: exchange,
			price:    price,
		})
	}

	// Helper to capture metadata updates
	captureMetadataUpdate := func(key marketTypes.UpdateKey) {
		metadataUpdateMutex.Lock()
		defer metadataUpdateMutex.Unlock()
		metadataUpdates = append(metadataUpdates, key)
	}

	BeforeEach(func() {
		// Reset captured updates
		priceUpdates = []priceUpdate{}
		metadataUpdates = []marketTypes.UpdateKey{}

		// Create extension with callbacks
		extension = extensions.NewOrderBookExtension(
			capturePriceUpdate,
			captureMetadataUpdate,
		)

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
		Context("when adding a new order book", func() {
			It("should store the order book correctly", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
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

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				retrieved := extension.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(retrieved.Bids).To(HaveLen(2))
				Expect(retrieved.Asks).To(HaveLen(2))
				Expect(retrieved.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(retrieved.Asks[0].Price).To(Equal(numerical.NewFromFloat(50010)))
			})

			It("should trigger price update callback with mid-price", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				Expect(priceUpdates).To(HaveLen(1))
				update := priceUpdates[0]
				Expect(update.pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(update.exchange).To(Equal(connector.ExchangeName("hyperliquid")))

				expectedMidPrice := numerical.NewFromFloat(50005) // (50000 + 50010) / 2
				Expect(update.price.Price.IntPart()).To(Equal(expectedMidPrice.IntPart()))
				Expect(update.price.BidPrice).To(Equal(numerical.NewFromFloat(50000)))
				Expect(update.price.AskPrice).To(Equal(numerical.NewFromFloat(50010)))
			})

			It("should trigger metadata update callback", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				Expect(metadataUpdates).To(HaveLen(1))
				update := metadataUpdates[0]
				Expect(update.DataType).To(Equal(marketTypes.DataKeyOrderBooks))
				Expect(update.Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(update.Exchange).To(Equal(connector.ExchangeName("hyperliquid")))
			})

			It("should handle multiple exchanges for the same pair", func() {
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

				extension.UpdateOrderBook(btc, "hyperliquid", orderBookHyper)
				extension.UpdateOrderBook(btc, "bybit", orderBookBybit)

				hyperBook := extension.GetOrderBook(btc, "hyperliquid")
				bybitBook := extension.GetOrderBook(btc, "bybit")

				Expect(hyperBook).NotTo(BeNil())
				Expect(bybitBook).NotTo(BeNil())
				Expect(hyperBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50000)))
				Expect(bybitBook.Bids[0].Price).To(Equal(numerical.NewFromFloat(50100)))
			})

			It("should handle multiple pairs", func() {
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

				extension.UpdateOrderBook(btc, "hyperliquid", orderBookBTC)
				extension.UpdateOrderBook(eth, "hyperliquid", orderBookETH)

				btcBook := extension.GetOrderBook(btc, "hyperliquid")
				ethBook := extension.GetOrderBook(eth, "hyperliquid")

				Expect(btcBook).NotTo(BeNil())
				Expect(ethBook).NotTo(BeNil())
				Expect(btcBook.Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(ethBook.Pair.Symbol()).To(Equal(eth.Symbol()))
			})
		})

		Context("when updating an existing order book", func() {
			It("should replace the order book with new data", func() {
				now := time.Now()

				orderBook1 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				orderBook2 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now.Add(time.Second),
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50005), Quantity: numerical.NewFromFloat(2.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50015), Quantity: numerical.NewFromFloat(1.5)},
					},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook1)
				extension.UpdateOrderBook(btc, "hyperliquid", orderBook2)

				retrieved := extension.GetOrderBook(btc, "hyperliquid")

				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Bids[0].Price).To(Equal(numerical.NewFromFloat(50005)))
				Expect(retrieved.Asks[0].Price).To(Equal(numerical.NewFromFloat(50015)))
				Expect(retrieved.Timestamp).To(Equal(now.Add(time.Second)))
			})

			It("should trigger callbacks on each update", func() {
				now := time.Now()

				orderBook1 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				orderBook2 := connector.OrderBook{
					Pair:      btc,
					Timestamp: now.Add(time.Second),
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50005), Quantity: numerical.NewFromFloat(2.0)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50015), Quantity: numerical.NewFromFloat(1.5)},
					},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook1)
				extension.UpdateOrderBook(btc, "hyperliquid", orderBook2)

				Expect(priceUpdates).To(HaveLen(2))
				Expect(metadataUpdates).To(HaveLen(2))
			})
		})

		Context("when dealing with edge cases", func() {
			It("should handle empty bids and asks without triggering price callback", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{},
					Asks:      []connector.PriceLevel{},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				retrieved := extension.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Bids).To(HaveLen(0))
				Expect(retrieved.Asks).To(HaveLen(0))

				// Should not trigger price update (no bids/asks to calculate mid-price)
				Expect(priceUpdates).To(HaveLen(0))

				// Should still trigger metadata update
				Expect(metadataUpdates).To(HaveLen(1))
			})

			It("should handle only bids without triggering price callback", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				// Should not trigger price update (need both bids and asks)
				Expect(priceUpdates).To(HaveLen(0))
			})

			It("should handle only asks without triggering price callback", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids:      []connector.PriceLevel{},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				// Should not trigger price update (need both bids and asks)
				Expect(priceUpdates).To(HaveLen(0))
			})
		})

		Context("with nil callbacks", func() {
			BeforeEach(func() {
				// Create extension without callbacks
				extension = extensions.NewOrderBookExtension(nil, nil)
			})

			It("should work without callbacks", func() {
				now := time.Now()
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}

				// Should not panic
				Expect(func() {
					extension.UpdateOrderBook(btc, "hyperliquid", orderBook)
				}).NotTo(Panic())

				retrieved := extension.GetOrderBook(btc, "hyperliquid")
				Expect(retrieved).NotTo(BeNil())
			})
		})

		Context("concurrent access", func() {
			It("should handle concurrent updates without data loss", func() {
				now := time.Now()
				var wg sync.WaitGroup

				// Spawn 10 goroutines updating different exchanges
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func(idx int64) {
						defer wg.Done()
						defer GinkgoRecover()

						exchangeName := connector.ExchangeName("exchange-" + string(rune('0'+idx)))
						orderBook := connector.OrderBook{
							Pair:      btc,
							Timestamp: now,
							Bids: []connector.PriceLevel{
								{Price: numerical.NewFromInt(50000 + idx), Quantity: numerical.NewFromFloat(1.5)},
							},
							Asks: []connector.PriceLevel{
								{Price: numerical.NewFromInt(50010 + idx), Quantity: numerical.NewFromFloat(1.0)},
							},
						}
						extension.UpdateOrderBook(btc, exchangeName, orderBook)
					}(int64(i))
				}

				wg.Wait()

				// Verify all updates were stored
				books := extension.GetOrderBooks(btc)
				Expect(books).To(HaveLen(10))

				// Verify each exchange has correct data
				for i := 0; i < 10; i++ {
					exchangeName := connector.ExchangeName("exchange-" + string(rune('0'+i)))
					book := extension.GetOrderBook(btc, exchangeName)
					Expect(book).NotTo(BeNil())
					Expect(book.Bids[0].Price).To(Equal(numerical.NewFromInt(50000 + int64(i))))
				}
			})

			It("should handle concurrent reads and writes", func() {
				now := time.Now()
				var wg sync.WaitGroup

				// Initial order book
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}
				extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

				// Spawn readers and writers
				for i := 0; i < 20; i++ {
					wg.Add(1)
					if i%2 == 0 {
						// Writer
						go func(idx int64) {
							defer wg.Done()
							defer GinkgoRecover()

							book := connector.OrderBook{
								Pair:      btc,
								Timestamp: now,
								Bids: []connector.PriceLevel{
									{Price: numerical.NewFromInt(50000 + idx), Quantity: numerical.NewFromFloat(1.5)},
								},
								Asks: []connector.PriceLevel{
									{Price: numerical.NewFromInt(50010 + idx), Quantity: numerical.NewFromFloat(1.0)},
								},
							}
							extension.UpdateOrderBook(btc, "hyperliquid", book)
						}(int64(i))
					} else {
						// Reader
						go func() {
							defer wg.Done()
							defer GinkgoRecover()

							book := extension.GetOrderBook(btc, "hyperliquid")
							Expect(book).NotTo(BeNil())
						}()
					}
				}

				wg.Wait()
			})
		})
	})

	Describe("GetOrderBook", func() {
		It("should return nil for non-existent pair", func() {
			book := extension.GetOrderBook(btc, "hyperliquid")
			Expect(book).To(BeNil())
		})

		It("should return nil for non-existent exchange", func() {
			now := time.Now()
			orderBook := connector.OrderBook{
				Pair:      btc,
				Timestamp: now,
				Bids: []connector.PriceLevel{
					{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
				},
				Asks: []connector.PriceLevel{
					{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
				},
			}
			extension.UpdateOrderBook(btc, "hyperliquid", orderBook)

			book := extension.GetOrderBook(btc, "bybit")
			Expect(book).To(BeNil())
		})
	})

	Describe("GetOrderBooks", func() {
		It("should return empty map for non-existent pair", func() {
			books := extension.GetOrderBooks(btc)
			Expect(books).To(BeEmpty())
		})

		It("should return all exchanges for a pair", func() {
			now := time.Now()

			exchanges := []connector.ExchangeName{"hyperliquid", "bybit", "binance"}
			for _, exchange := range exchanges {
				orderBook := connector.OrderBook{
					Pair:      btc,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}
				extension.UpdateOrderBook(btc, exchange, orderBook)
			}

			books := extension.GetOrderBooks(btc)
			Expect(books).To(HaveLen(3))
			Expect(books).To(HaveKey(connector.ExchangeName("hyperliquid")))
			Expect(books).To(HaveKey(connector.ExchangeName("bybit")))
			Expect(books).To(HaveKey(connector.ExchangeName("binance")))
		})
	})

	Describe("GetAllPairsWithOrderBooks", func() {
		It("should return empty slice when no order books", func() {
			pairs := extension.GetAllPairsWithOrderBooks()
			Expect(pairs).To(BeEmpty())
		})

		It("should return all pairs with order books", func() {
			now := time.Now()

			// Add order books for multiple pairs
			pairs := []portfolio.Pair{btc, eth}
			for _, pair := range pairs {
				orderBook := connector.OrderBook{
					Pair:      pair,
					Timestamp: now,
					Bids: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50000), Quantity: numerical.NewFromFloat(1.5)},
					},
					Asks: []connector.PriceLevel{
						{Price: numerical.NewFromFloat(50010), Quantity: numerical.NewFromFloat(1.0)},
					},
				}
				extension.UpdateOrderBook(pair, "hyperliquid", orderBook)
			}

			allPairs := extension.GetAllPairsWithOrderBooks()
			Expect(allPairs).To(HaveLen(2))

			// Convert to symbols for easier comparison
			symbols := make([]string, len(allPairs))
			for i, p := range allPairs {
				symbols[i] = p.Symbol()
			}
			Expect(symbols).To(ContainElements(btc.Symbol(), eth.Symbol()))
		})
	})
})
