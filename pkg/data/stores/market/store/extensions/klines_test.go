package extensions_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/market/store/extensions"
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

var _ = Describe("Kline Store", func() {
	var (
		extension marketTypes.KlineStoreExtension
		btc       portfolio.Pair
		eth       portfolio.Pair
	)

	BeforeEach(func() {
		// Create extension directly - no dependencies needed
		extension = extensions.NewKlineExtension()

		btc = portfolio.NewPair(
			portfolio.NewAsset("BTC"),
			portfolio.NewAsset("USD"),
		)

		eth = portfolio.NewPair(
			portfolio.NewAsset("ETH"),
			portfolio.NewAsset("USD"),
		)
	})

	Describe("UpdateKline", func() {
		Context("when adding a new kline", func() {
			It("should store the kline correctly", func() {
				now := time.Now().Truncate(time.Minute)
				kline := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					High:      50100,
					Low:       49900,
					Close:     50050,
					Volume:    100,
				}

				extension.UpdateKline(btc, "hyperliquid", kline)

				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				Expect(klines).To(HaveLen(1))
				Expect(klines[0].Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(klines[0].Open).To(Equal(50000.0))
				Expect(klines[0].Close).To(Equal(50050.0))
				Expect(klines[0].High).To(Equal(50100.0))
				Expect(klines[0].Low).To(Equal(49900.0))
			})

			It("should handle multiple intervals for the same pair and exchange", func() {
				now := time.Now().Truncate(time.Minute)

				kline1m := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}

				kline5m := connector.Kline{
					Pair:      btc,
					Interval:  "5m",
					OpenTime:  now,
					CloseTime: now.Add(5 * time.Minute),
					Open:      50000,
					Close:     50200,
				}

				extension.UpdateKline(btc, "hyperliquid", kline1m)
				extension.UpdateKline(btc, "hyperliquid", kline5m)

				klines1m := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				klines5m := extension.GetKlines(btc, "hyperliquid", "5m", 0)

				Expect(klines1m).To(HaveLen(1))
				Expect(klines5m).To(HaveLen(1))
				Expect(klines1m[0].Interval).To(Equal("1m"))
				Expect(klines5m[0].Interval).To(Equal("5m"))
			})

			It("should handle multiple exchanges for the same pair", func() {
				now := time.Now().Truncate(time.Minute)

				klineHyper := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}

				klineBybit := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50100,
					Close:     50150,
				}

				extension.UpdateKline(btc, "hyperliquid", klineHyper)
				extension.UpdateKline(btc, "bybit", klineBybit)

				klinesHyper := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				klinesBybit := extension.GetKlines(btc, "bybit", "1m", 0)

				Expect(klinesHyper).To(HaveLen(1))
				Expect(klinesBybit).To(HaveLen(1))
				Expect(klinesHyper[0].Open).To(Equal(float64(50000)))
				Expect(klinesBybit[0].Open).To(Equal(float64(50100)))
			})

			It("should handle multiple pairs", func() {
				now := time.Now().Truncate(time.Minute)

				klineBTC := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}

				klineETH := connector.Kline{
					Pair:      eth,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      3000,
					Close:     3010,
				}

				extension.UpdateKline(btc, "hyperliquid", klineBTC)
				extension.UpdateKline(eth, "hyperliquid", klineETH)

				klinesBTC := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				klinesETH := extension.GetKlines(eth, "hyperliquid", "1m", 0)

				Expect(klinesBTC).To(HaveLen(1))
				Expect(klinesETH).To(HaveLen(1))
				Expect(klinesBTC[0].Pair.Symbol()).To(Equal(btc.Symbol()))
				Expect(klinesETH[0].Pair.Symbol()).To(Equal(eth.Symbol()))
			})
		})

		Context("when updating an existing kline", func() {
			It("should replace the kline with the same open time", func() {
				now := time.Now().Truncate(time.Minute)

				// First kline
				kline1 := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
					High:      50100,
					Low:       49900,
				}

				// Updated kline with same open time
				kline2 := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now, // Same open time
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50150, // Different close
					High:      50200, // Different high
					Low:       49800, // Different low
				}

				extension.UpdateKline(btc, "hyperliquid", kline1)
				extension.UpdateKline(btc, "hyperliquid", kline2)

				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)

				Expect(klines).To(HaveLen(1), "Should only have one kline after update")
				Expect(klines[0].Close).To(Equal(50150.0), "Close should be updated")
				Expect(klines[0].High).To(Equal(50200.0), "High should be updated")
				Expect(klines[0].Low).To(Equal(49800.0), "Low should be updated")
			})

			It("should append kline with different open time", func() {
				now := time.Now().Truncate(time.Minute)

				kline1 := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}

				kline2 := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now.Add(time.Minute), // Different open time
					CloseTime: now.Add(2 * time.Minute),
					Open:      50050,
					Close:     50100,
				}

				extension.UpdateKline(btc, "hyperliquid", kline1)
				extension.UpdateKline(btc, "hyperliquid", kline2)

				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)

				Expect(klines).To(HaveLen(2), "Should have two klines")
				Expect(klines[0].OpenTime).To(Equal(now))
				Expect(klines[1].OpenTime).To(Equal(now.Add(time.Minute)))
			})
		})

		Context("when adding multiple sequential klines", func() {
			It("should store them in order", func() {
				now := time.Now().Truncate(time.Minute)

				// Add 5 sequential klines
				for i := 0; i < 5; i++ {
					kline := connector.Kline{
						Pair:      btc,
						Interval:  "1m",
						OpenTime:  now.Add(time.Duration(i) * time.Minute),
						CloseTime: now.Add(time.Duration(i+1) * time.Minute),
						Open:      50000 + float64(i*10),
						Close:     50010 + float64(i*10),
					}
					extension.UpdateKline(btc, "hyperliquid", kline)
				}

				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				Expect(klines).To(HaveLen(5))

				// Verify they're stored correctly
				for i := 0; i < 5; i++ {
					expectedOpen := 50000.0 + float64(i*10)
					Expect(klines[i].Open).To(Equal(expectedOpen))
				}
			})
		})

		Context("concurrent access", func() {
			It("should handle concurrent updates without data loss", func() {
				now := time.Now().Truncate(time.Minute)
				var wg sync.WaitGroup

				// Spawn 10 goroutines updating different time periods
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func(idx int) {
						defer wg.Done()
						defer GinkgoRecover()

						kline := connector.Kline{
							Pair:      btc,
							Interval:  "1m",
							OpenTime:  now.Add(time.Duration(idx) * time.Minute),
							CloseTime: now.Add(time.Duration(idx+1) * time.Minute),
							Open:      50000 + float64(idx*10),
							Close:     50010 + float64(idx*10),
						}
						extension.UpdateKline(btc, "hyperliquid", kline)
					}(i)
				}

				wg.Wait()

				// Verify all updates were stored
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				Expect(klines).To(HaveLen(10))
			})

			It("should handle concurrent updates to different exchanges", func() {
				now := time.Now().Truncate(time.Minute)
				var wg sync.WaitGroup

				// Spawn goroutines for different exchanges
				exchanges := []connector.ExchangeName{"hyperliquid", "bybit", "binance", "coinbase", "kraken"}
				for _, exchange := range exchanges {
					wg.Add(1)
					go func(ex connector.ExchangeName) {
						defer wg.Done()
						defer GinkgoRecover()

						for i := 0; i < 5; i++ {
							kline := connector.Kline{
								Pair:      btc,
								Interval:  "1m",
								OpenTime:  now.Add(time.Duration(i) * time.Minute),
								CloseTime: now.Add(time.Duration(i+1) * time.Minute),
								Open:      50000 + float64(i*10),
								Close:     50010 + float64(i*10),
							}
							extension.UpdateKline(btc, ex, kline)
						}
					}(exchange)
				}

				wg.Wait()

				// Verify all exchanges have data
				for _, exchange := range exchanges {
					klines := extension.GetKlines(btc, exchange, "1m", 0)
					Expect(klines).To(HaveLen(5))
				}
			})

			It("should handle concurrent reads and writes", func() {
				now := time.Now().Truncate(time.Minute)
				var wg sync.WaitGroup

				// Pre-populate with some data
				for i := 0; i < 5; i++ {
					kline := connector.Kline{
						Pair:      btc,
						Interval:  "1m",
						OpenTime:  now.Add(time.Duration(i) * time.Minute),
						CloseTime: now.Add(time.Duration(i+1) * time.Minute),
						Open:      50000 + float64(i*10),
						Close:     50010 + float64(i*10),
					}
					extension.UpdateKline(btc, "hyperliquid", kline)
				}

				// Spawn readers and writers
				for i := 0; i < 20; i++ {
					wg.Add(1)
					if i%2 == 0 {
						// Writer
						go func(idx int) {
							defer wg.Done()
							defer GinkgoRecover()

							kline := connector.Kline{
								Pair:      btc,
								Interval:  "1m",
								OpenTime:  now.Add(time.Duration(5+idx) * time.Minute),
								CloseTime: now.Add(time.Duration(6+idx) * time.Minute),
								Open:      50000 + float64((5+idx)*10),
								Close:     50010 + float64((5+idx)*10),
							}
							extension.UpdateKline(btc, "hyperliquid", kline)
						}(i)
					} else {
						// Reader
						go func() {
							defer wg.Done()
							defer GinkgoRecover()

							klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
							Expect(klines).NotTo(BeEmpty())
						}()
					}
				}

				wg.Wait()
			})
		})
	})

	Describe("GetKlines", func() {
		Context("when querying non-existent data", func() {
			It("should return empty slice for non-existent pair", func() {
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				Expect(klines).To(BeEmpty())
			})

			It("should return empty slice for non-existent exchange", func() {
				now := time.Now().Truncate(time.Minute)
				kline := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}
				extension.UpdateKline(btc, "hyperliquid", kline)

				klines := extension.GetKlines(btc, "bybit", "1m", 0)
				Expect(klines).To(BeEmpty())
			})

			It("should return empty slice for non-existent interval", func() {
				now := time.Now().Truncate(time.Minute)
				kline := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}
				extension.UpdateKline(btc, "hyperliquid", kline)

				klines := extension.GetKlines(btc, "hyperliquid", "5m", 0)
				Expect(klines).To(BeEmpty())
			})
		})

		Context("with limit parameter", func() {
			BeforeEach(func() {
				now := time.Now().Truncate(time.Minute)
				// Add 10 klines
				for i := 0; i < 10; i++ {
					kline := connector.Kline{
						Pair:      btc,
						Interval:  "1m",
						OpenTime:  now.Add(time.Duration(i) * time.Minute),
						CloseTime: now.Add(time.Duration(i+1) * time.Minute),
						Open:      50000 + float64(i*10),
						Close:     50010 + float64(i*10),
					}
					extension.UpdateKline(btc, "hyperliquid", kline)
				}
			})

			It("should return all klines when limit is 0", func() {
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
				Expect(klines).To(HaveLen(10))
			})

			It("should return limited klines from the end", func() {
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 5)
				Expect(klines).To(HaveLen(5))

				// Should get the last 5 klines
				Expect(klines[0].Open).To(Equal(50050.0)) // 6th kline
				Expect(klines[4].Open).To(Equal(50090.0)) // 10th kline
			})

			It("should return all klines when limit exceeds available", func() {
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 100)
				Expect(klines).To(HaveLen(10))
			})

			It("should handle limit of 1", func() {
				klines := extension.GetKlines(btc, "hyperliquid", "1m", 1)
				Expect(klines).To(HaveLen(1))
				Expect(klines[0].Open).To(Equal(50090.0)) // Last kline
			})
		})
	})

	Describe("GetKlinesSince", func() {
		var now time.Time

		BeforeEach(func() {
			now = time.Now().Truncate(time.Minute)
			// Add 10 klines spanning 10 minutes
			for i := 0; i < 10; i++ {
				kline := connector.Kline{
					Pair:      btc,
					Interval:  "1m",
					OpenTime:  now.Add(time.Duration(i) * time.Minute),
					CloseTime: now.Add(time.Duration(i+1) * time.Minute),
					Open:      50000 + float64(i*10),
					Close:     50010 + float64(i*10),
				}
				extension.UpdateKline(btc, "hyperliquid", kline)
			}
		})

		It("should return all klines when since is before first kline", func() {
			since := now.Add(-5 * time.Minute)
			klines := extension.GetKlinesSince(btc, "hyperliquid", "1m", since)
			Expect(klines).To(HaveLen(10))
		})

		It("should return klines at or after the since time", func() {
			since := now.Add(5 * time.Minute)
			klines := extension.GetKlinesSince(btc, "hyperliquid", "1m", since)

			// Should include klines from minute 5 onwards (5, 6, 7, 8, 9)
			Expect(klines).To(HaveLen(5))
			Expect(klines[0].Open).To(Equal(50050.0)) // 6th kline
		})

		It("should include kline with exact since time", func() {
			since := now.Add(5 * time.Minute)
			klines := extension.GetKlinesSince(btc, "hyperliquid", "1m", since)

			// Should include the kline that starts at exactly 'since'
			Expect(klines[0].OpenTime).To(Equal(since))
		})

		It("should return empty slice when since is after all klines", func() {
			since := now.Add(20 * time.Minute)
			klines := extension.GetKlinesSince(btc, "hyperliquid", "1m", since)
			Expect(klines).To(BeEmpty())
		})

		It("should return empty slice for non-existent pair", func() {
			since := now
			klines := extension.GetKlinesSince(eth, "hyperliquid", "1m", since)
			Expect(klines).To(BeEmpty())
		})

		It("should return empty slice for non-existent exchange", func() {
			since := now
			klines := extension.GetKlinesSince(btc, "bybit", "1m", since)
			Expect(klines).To(BeEmpty())
		})

		It("should return empty slice for non-existent interval", func() {
			since := now
			klines := extension.GetKlinesSince(btc, "hyperliquid", "5m", since)
			Expect(klines).To(BeEmpty())
		})

		It("should handle multiple intervals independently", func() {
			// Add some 5m klines
			for i := 0; i < 5; i++ {
				kline := connector.Kline{
					Pair:      btc,
					Interval:  "5m",
					OpenTime:  now.Add(time.Duration(i*5) * time.Minute),
					CloseTime: now.Add(time.Duration((i+1)*5) * time.Minute),
					Open:      50000 + float64(i*50),
					Close:     50050 + float64(i*50),
				}
				extension.UpdateKline(btc, "hyperliquid", kline)
			}

			since := now.Add(10 * time.Minute)

			// Get 1m klines
			klines1m := extension.GetKlinesSince(btc, "hyperliquid", "1m", since)
			Expect(klines1m).To(BeEmpty())

			// Get 5m klines
			klines5m := extension.GetKlinesSince(btc, "hyperliquid", "5m", since)
			Expect(klines5m).To(HaveLen(3)) // Should have klines at 10, 15, 20 minutes
		})
	})

	Describe("Edge Cases", func() {
		It("should handle zero values correctly", func() {
			now := time.Now().Truncate(time.Minute)
			kline := connector.Kline{
				Pair:      btc,
				Interval:  "1m",
				OpenTime:  now,
				CloseTime: now.Add(time.Minute),
				Open:      0,
				High:      0,
				Low:       0,
				Close:     0,
				Volume:    0,
			}

			extension.UpdateKline(btc, "hyperliquid", kline)

			klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
			Expect(klines).To(HaveLen(1))
			Expect(klines[0].Open).To(Equal(0.0))
			Expect(klines[0].Close).To(Equal(0.0))
		})

		It("should handle very large numbers", func() {
			now := time.Now().Truncate(time.Minute)
			kline := connector.Kline{
				Pair:      btc,
				Interval:  "1m",
				OpenTime:  now,
				CloseTime: now.Add(time.Minute),
				Open:      999999999999.99,
				High:      999999999999.99,
				Low:       999999999999.99,
				Close:     999999999999.99,
				Volume:    999999999999.99,
			}

			extension.UpdateKline(btc, "hyperliquid", kline)

			klines := extension.GetKlines(btc, "hyperliquid", "1m", 0)
			Expect(klines).To(HaveLen(1))
			Expect(klines[0].Open).To(Equal(999999999999.99))
		})

		It("should handle special interval strings", func() {
			now := time.Now().Truncate(time.Minute)

			intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d", "1w"}
			for _, interval := range intervals {
				kline := connector.Kline{
					Pair:      btc,
					Interval:  interval,
					OpenTime:  now,
					CloseTime: now.Add(time.Minute),
					Open:      50000,
					Close:     50050,
				}
				extension.UpdateKline(btc, "hyperliquid", kline)
			}

			// Verify all intervals are stored independently
			for _, interval := range intervals {
				klines := extension.GetKlines(btc, "hyperliquid", interval, 0)
				Expect(klines).To(HaveLen(1), "Should have kline for interval: "+interval)
			}
		})
	})
})
