package trade_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/data/stores/activity/trade"
	activityTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("Trade Store - Add Trades", func() {
	var (
		store   activityTypes.Trades
		btcPair portfolio.Pair
		ethPair portfolio.Pair
		solPair portfolio.Pair
	)

	BeforeEach(func() {
		store = trade.NewStore()
		btcPair = portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USDT"))
		ethPair = portfolio.NewPair(portfolio.NewAsset("ETH"), portfolio.NewAsset("USDT"))
		solPair = portfolio.NewPair(portfolio.NewAsset("SOL"), portfolio.NewAsset("USDT"))
	})

	Describe("AddTrade", func() {
		Context("when adding a new trade", func() {
			It("should store the trade correctly", func() {
				t := connector.Trade{
					ID:        "trade-1",
					OrderID:   "order-1",
					Pair:      btcPair,
					Exchange:  "hyperliquid",
					Price:     numerical.NewFromFloat(50000),
					Quantity:  numerical.NewFromFloat(1.0),
					Side:      connector.OrderSideBuy,
					IsMaker:   false,
					Fee:       numerical.NewFromFloat(10),
					Timestamp: time.Now(),
				}

				store.AddTrade(t)

				Expect(store.GetTradeCount()).To(Equal(1))
				retrieved := store.GetTradeByID("trade-1")
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Pair.Symbol()).To(Equal("BTC-USDT"))
				Expect(retrieved.Price.Equal(numerical.NewFromFloat(50000))).To(BeTrue())
			})

			It("should handle multiple trades", func() {
				for i := 1; i <= 5; i++ {
					t := connector.Trade{
						ID:       "trade-" + string(rune('0'+i)),
						Pair:     btcPair,
						Exchange: "hyperliquid",
					}
					store.AddTrade(t)
				}

				Expect(store.GetTradeCount()).To(Equal(5))
			})

			It("should skip duplicate trades", func() {
				t := connector.Trade{
					ID:   "trade-1",
					Pair: btcPair,
				}

				store.AddTrade(t)
				store.AddTrade(t) // Same ID

				Expect(store.GetTradeCount()).To(Equal(1))
			})

			It("should preserve insertion order", func() {
				now := time.Now()
				trades := []connector.Trade{
					{ID: "t1", Pair: btcPair, Timestamp: now},
					{ID: "t2", Pair: ethPair, Timestamp: now.Add(time.Second)},
					{ID: "t3", Pair: solPair, Timestamp: now.Add(2 * time.Second)},
				}

				for _, t := range trades {
					store.AddTrade(t)
				}

				all := store.GetAllTrades()
				Expect(all).To(HaveLen(3))
				Expect(all[0].ID).To(Equal("t1"))
				Expect(all[1].ID).To(Equal("t2"))
				Expect(all[2].ID).To(Equal("t3"))
			})
		})
	})

	Describe("AddTrades", func() {
		Context("when adding multiple trades at once", func() {
			It("should store all trades correctly", func() {
				trades := []connector.Trade{
					{ID: "t1", Pair: btcPair, Exchange: "hyperliquid"},
					{ID: "t2", Pair: ethPair, Exchange: "hyperliquid"},
					{ID: "t3", Pair: solPair, Exchange: "bybit"},
				}

				store.AddTrades(trades)

				Expect(store.GetTradeCount()).To(Equal(3))
				Expect(store.TradeExists("t1")).To(BeTrue())
				Expect(store.TradeExists("t2")).To(BeTrue())
				Expect(store.TradeExists("t3")).To(BeTrue())
			})

			It("should skip duplicate trades in batch against existing trades", func() {
				// Add initial trade
				store.AddTrade(connector.Trade{ID: "t1", Pair: btcPair})

				// Add batch with duplicate of existing trade
				trades := []connector.Trade{
					{ID: "t1", Pair: btcPair}, // Duplicate of existing
					{ID: "t2", Pair: ethPair},
				}

				store.AddTrades(trades)

				Expect(store.GetTradeCount()).To(Equal(2))
			})

			It("should handle empty slice", func() {
				store.AddTrades([]connector.Trade{})

				Expect(store.GetTradeCount()).To(Equal(0))
			})

			It("should handle all duplicates", func() {
				store.AddTrade(connector.Trade{ID: "t1", Pair: btcPair})
				store.AddTrade(connector.Trade{ID: "t2", Pair: ethPair})

				// Add batch with all existing trades
				trades := []connector.Trade{
					{ID: "t1", Pair: btcPair},
					{ID: "t2", Pair: ethPair},
				}

				store.AddTrades(trades)

				Expect(store.GetTradeCount()).To(Equal(2))
			})
		})
	})

	Describe("Clear", func() {
		Context("when clearing the store", func() {
			It("should remove all trades", func() {
				store.AddTrade(connector.Trade{ID: "t1", Pair: btcPair})
				store.AddTrade(connector.Trade{ID: "t2", Pair: ethPair})

				store.Clear()

				Expect(store.GetTradeCount()).To(Equal(0))
				Expect(store.GetAllTrades()).To(BeEmpty())
			})

			It("should reset ID lookup", func() {
				store.AddTrade(connector.Trade{ID: "t1", Pair: btcPair})

				store.Clear()

				Expect(store.TradeExists("t1")).To(BeFalse())
				Expect(store.GetTradeByID("t1")).To(BeNil())
			})
		})
	})

	Describe("Concurrent access", func() {
		Context("when multiple goroutines add trades", func() {
			It("should handle concurrent writes safely", func() {
				done := make(chan bool)
				iterations := 50

				// Writer 1
				go func() {
					for i := 0; i < iterations; i++ {
						store.AddTrade(connector.Trade{
							ID:   "btc-" + string(rune(i)),
							Pair: btcPair,
						})
					}
					done <- true
				}()

				// Writer 2
				go func() {
					for i := 0; i < iterations; i++ {
						store.AddTrade(connector.Trade{
							ID:   "eth-" + string(rune(i)),
							Pair: ethPair,
						})
					}
					done <- true
				}()

				// Reader
				go func() {
					for i := 0; i < iterations; i++ {
						_ = store.GetAllTrades()
						_ = store.GetTradeCount()
					}
					done <- true
				}()

				<-done
				<-done
				<-done

				Expect(store.GetTradeCount()).To(Equal(iterations * 2))
			})
		})
	})
})
