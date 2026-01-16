package activity_test

import (
	"context"
	"time"

	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	storeActivity "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("PNL", func() {
	var (
		app            *fxtest.App
		pnl            kronosActivity.PNL
		positionStore  storeActivity.Positions
		tradesStore    storeActivity.Trades
		marketRegistry marketTypes.MarketRegistry
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&pnl,
				&positionStore,
				&tradesStore,
				&marketRegistry,
			),
			fx.NopLogger,
		)

		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	// Helper to set asset price in spot market
	setAssetPrice := func(asset portfolio.Asset, price float64) {
		spotStore := marketRegistry.Get(marketTypes.MarketTypeSpot)
		spotStore.UpdateAssetPrice(asset, "binance", connector.Price{
			Price:     numerical.NewFromFloat(price),
			Timestamp: time.Now(),
		})
	}

	Describe("GetFeesByStrategy", func() {
		It("should sum fees for a strategy", func() {
			ctx := context.Background()
			strategyName := strategy.StrategyName("test-strategy")
			ctx = strategy.WithStrategyName(ctx, strategyName)

			// Add trades to position store
			positionStore.AddTradeToStrategy(strategyName, connector.Trade{ID: "t1", Symbol: "BTC", Fee: numerical.NewFromFloat(10)})
			positionStore.AddTradeToStrategy(strategyName, connector.Trade{ID: "t2", Symbol: "BTC", Fee: numerical.NewFromFloat(15)})
			positionStore.AddTradeToStrategy(strategyName, connector.Trade{ID: "t3", Symbol: "ETH", Fee: numerical.NewFromFloat(5)})

			result := pnl.GetFeesByStrategy(ctx, strategyName)

			Expect(result.Equal(numerical.NewFromFloat(30))).To(BeTrue())
		})

		It("should return zero when no trades", func() {
			ctx := context.Background()
			strategyName := strategy.StrategyName("empty-strategy")

			result := pnl.GetFeesByStrategy(ctx, strategyName)

			Expect(result.IsZero()).To(BeTrue())
		})
	})

	Describe("GetTotalFees", func() {
		It("should sum fees across all trades", func() {
			ctx := context.Background()

			// Add trades to trades store
			tradesStore.AddTrade(connector.Trade{ID: "t1", Fee: numerical.NewFromFloat(10)})
			tradesStore.AddTrade(connector.Trade{ID: "t2", Fee: numerical.NewFromFloat(20)})
			tradesStore.AddTrade(connector.Trade{ID: "t3", Fee: numerical.NewFromFloat(30)})

			result := pnl.GetTotalFees(ctx)

			Expect(result.Equal(numerical.NewFromFloat(60))).To(BeTrue())
		})
	})

	Describe("GetRealizedPNL", func() {
		Context("long position closed for profit", func() {
			It("should calculate profit correctly", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Buy 1 BTC at 50000, sell 1 BTC at 55000 = profit of 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:        "t1",
					Symbol:    "BTC",
					Side:      connector.OrderSideBuy,
					Quantity:  numerical.NewFromFloat(1),
					Price:     numerical.NewFromFloat(50000),
					Fee:       numerical.NewFromFloat(50),
					Timestamp: time.Now(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:        "t2",
					Symbol:    "BTC",
					Side:      connector.OrderSideSell,
					Quantity:  numerical.NewFromFloat(1),
					Price:     numerical.NewFromFloat(55000),
					Fee:       numerical.NewFromFloat(55),
					Timestamp: time.Now(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				// PNL = 5000 - fees (50 + 55) = 4895
				Expect(result.Equal(numerical.NewFromFloat(4895))).To(BeTrue())
			})
		})

		Context("long position closed for loss", func() {
			It("should calculate loss correctly", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)
				
				// Buy 1 BTC at 50000, sell 1 BTC at 45000 = loss of 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(45000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				Expect(result.Equal(numerical.NewFromFloat(-5000))).To(BeTrue())
			})
		})

		Context("short position closed for profit", func() {
			It("should calculate profit correctly", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Sell 1 BTC at 50000, buy 1 BTC at 45000 = profit of 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(45000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})

		Context("short position closed for loss", func() {
			It("should calculate loss correctly", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Sell 1 BTC at 50000, buy 1 BTC at 55000 = loss of 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(55000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				Expect(result.Equal(numerical.NewFromFloat(-5000))).To(BeTrue())
			})
		})

		Context("partial close", func() {
			It("should calculate realized PNL for partial close only", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Buy 2 BTC at 50000, sell 1 BTC at 55000 = profit of 5000 on closed portion
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(55000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				// Only 1 BTC closed: (55000 - 50000) * 1 = 5000
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})

		Context("average cost basis", func() {
			It("should use weighted average entry price", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Buy 1 BTC at 50000, buy 1 BTC at 52000, sell 2 BTC at 54000
				// Avg entry = (50000 + 52000) / 2 = 51000
				// PNL = (54000 - 51000) * 2 = 6000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(52000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t3",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(54000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				Expect(result.Equal(numerical.NewFromFloat(6000))).To(BeTrue())
			})
		})

		Context("position flipping long to short", func() {
			It("should realize PNL on closed portion and track new short position", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Buy 1 BTC at 50000, sell 2 BTC at 55000
				// Closes long 1 BTC for profit of 5000, opens short 1 BTC at 55000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(55000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				// Only realized PNL from closing the long: (55000 - 50000) * 1 = 5000
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})

		Context("position flipping short to long", func() {
			It("should realize PNL on closed portion and track new long position", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Sell 1 BTC at 50000, buy 2 BTC at 45000
				// Closes short 1 BTC for profit of 5000, opens long 1 BTC at 45000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(45000),
					Fee:      numerical.Zero(),
				})

				result := pnl.GetRealizedPNL(ctx, strategyName)

				// Only realized PNL from closing the short: (50000 - 45000) * 1 = 5000
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})
	})

	Describe("GetUnrealizedPNL", func() {
		Context("open long position", func() {
			It("should calculate unrealized profit", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				btc := portfolio.NewAsset("BTC")

				// Buy 1 BTC at 50000, current price 55000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})

				setAssetPrice(btc, 55000)

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})

			It("should calculate unrealized loss", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				btc := portfolio.NewAsset("BTC")

				// Buy 1 BTC at 50000, current price 45000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})

				setAssetPrice(btc, 45000)

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Equal(numerical.NewFromFloat(-5000))).To(BeTrue())
			})
		})

		Context("open short position", func() {
			It("should calculate unrealized profit", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				btc := portfolio.NewAsset("BTC")

				// Sell 1 BTC at 50000, current price 45000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})

				setAssetPrice(btc, 45000)

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})

		Context("closed position", func() {
			It("should return zero unrealized PNL", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				// Buy 1 BTC, then sell 1 BTC - position is flat
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(55000),
					Fee:      numerical.Zero(),
				})

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.IsZero()).To(BeTrue())
			})
		})

		Context("flipped position unrealized PNL", func() {
			It("should use new entry price for flipped short position", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				btc := portfolio.NewAsset("BTC")

				// Buy 1 BTC at 50000, sell 2 BTC at 55000 (now short 1 at 55000)
				// Current price 60000 - short is losing 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(55000),
					Fee:      numerical.Zero(),
				})

				setAssetPrice(btc, 60000)

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				// Short 1 BTC at 55000, current price 60000: (55000 - 60000) * 1 = -5000
				Expect(result.Equal(numerical.NewFromFloat(-5000))).To(BeTrue())
			})

			It("should use new entry price for flipped long position", func() {
				ctx := context.Background()
				strategyName := strategy.StrategyName("test-strategy")
				ctx = strategy.WithStrategyName(ctx, strategyName)

				btc := portfolio.NewAsset("BTC")

				// Sell 1 BTC at 50000, buy 2 BTC at 45000 (now long 1 at 45000)
				// Current price 50000 - long is gaining 5000
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t1",
					Symbol:   "BTC",
					Side:     connector.OrderSideSell,
					Quantity: numerical.NewFromFloat(1),
					Price:    numerical.NewFromFloat(50000),
					Fee:      numerical.Zero(),
				})
				positionStore.AddTradeToStrategy(strategyName, connector.Trade{
					ID:       "t2",
					Symbol:   "BTC",
					Side:     connector.OrderSideBuy,
					Quantity: numerical.NewFromFloat(2),
					Price:    numerical.NewFromFloat(45000),
					Fee:      numerical.Zero(),
				})

				setAssetPrice(btc, 50000)

				result, err := pnl.GetUnrealizedPNL(ctx, strategyName)

				Expect(err).NotTo(HaveOccurred())
				// Long 1 BTC at 45000, current price 50000: (50000 - 45000) * 1 = 5000
				Expect(result.Equal(numerical.NewFromFloat(5000))).To(BeTrue())
			})
		})
	})

	Describe("GetTotalPNL", func() {
		It("should sum realized and unrealized PNL", func() {
			ctx := context.Background()
			btc := portfolio.NewAsset("BTC")
			eth := portfolio.NewAsset("ETH")

			// Closed BTC trade: profit of 5000
			tradesStore.AddTrade(connector.Trade{ID: "t1", Symbol: "BTC", Side: connector.OrderSideBuy, Quantity: numerical.NewFromFloat(1), Price: numerical.NewFromFloat(50000), Fee: numerical.Zero()})
			tradesStore.AddTrade(connector.Trade{ID: "t2", Symbol: "BTC", Side: connector.OrderSideSell, Quantity: numerical.NewFromFloat(1), Price: numerical.NewFromFloat(55000), Fee: numerical.Zero()})

			// Open ETH position
			tradesStore.AddTrade(connector.Trade{ID: "t3", Symbol: "ETH", Side: connector.OrderSideBuy, Quantity: numerical.NewFromFloat(10), Price: numerical.NewFromFloat(3000), Fee: numerical.Zero()})

			// Set prices
			setAssetPrice(btc, 55000)
			setAssetPrice(eth, 3300)

			result, err := pnl.GetTotalPNL(ctx)

			Expect(err).NotTo(HaveOccurred())
			// Realized: 5000, Unrealized: (3300-3000)*10 = 3000, Total: 8000
			Expect(result.Equal(numerical.NewFromFloat(8000))).To(BeTrue())
		})
	})
})
