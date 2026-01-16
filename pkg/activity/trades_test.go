package activity_test

import (
	"context"
	"time"

	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	storeActivity "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Trades", func() {
	var (
		app    *fxtest.App
		trades kronosActivity.Trades
		store  storeActivity.Trades
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&trades,
				&store,
			),
			fx.NopLogger,
		)

		app.RequireStart()
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Describe("GetAllTrades", func() {
		It("should return all trades from underlying store", func() {
			ctx := context.Background()

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Symbol: "BTC", Price: numerical.NewFromFloat(50000)})
			store.AddTrade(connector.Trade{ID: "trade-2", Symbol: "ETH", Price: numerical.NewFromFloat(3000)})

			result := trades.GetAllTrades(ctx)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("trade-1"))
			Expect(result[1].ID).To(Equal("trade-2"))
		})

		It("should return empty slice when no trades exist", func() {
			ctx := context.Background()

			result := trades.GetAllTrades(ctx)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByExchange", func() {
		It("should filter trades by exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("binance")

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Exchange: exchange, Symbol: "BTC"})
			store.AddTrade(connector.Trade{ID: "trade-2", Exchange: exchange, Symbol: "ETH"})
			store.AddTrade(connector.Trade{ID: "trade-3", Exchange: "coinbase", Symbol: "BTC"})

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Exchange).To(Equal(exchange))
			Expect(result[1].Exchange).To(Equal(exchange))
		})

		It("should return empty slice when no trades for exchange", func() {
			ctx := context.Background()

			exchange := connector.ExchangeName("unknown-exchange")

			result := trades.GetTradesByExchange(ctx, exchange)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesByAsset", func() {
		It("should filter trades by asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("BTC")

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1", Symbol: "BTC", Price: numerical.NewFromFloat(50000)})
			store.AddTrade(connector.Trade{ID: "trade-2", Symbol: "BTC", Price: numerical.NewFromFloat(51000)})
			store.AddTrade(connector.Trade{ID: "trade-3", Symbol: "ETH", Price: numerical.NewFromFloat(3000)})

			result := trades.GetTradesByAsset(ctx, asset)

			Expect(result).To(HaveLen(2))
			Expect(result[0].Symbol).To(Equal("BTC"))
			Expect(result[1].Symbol).To(Equal("BTC"))
		})

		It("should return empty slice when no trades for asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("UNKNOWN")

			result := trades.GetTradesByAsset(ctx, asset)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradesSince", func() {
		It("should filter trades by time", func() {
			ctx := context.Background()

			since := time.Now().Add(-1 * time.Hour)
			now := time.Now()

			// Add trades to store with different timestamps
			store.AddTrade(connector.Trade{ID: "trade-1", Timestamp: now})
			store.AddTrade(connector.Trade{ID: "trade-2", Timestamp: now.Add(-30 * time.Minute)})
			store.AddTrade(connector.Trade{ID: "trade-3", Timestamp: now.Add(-2 * time.Hour)}) // Before 'since'

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(HaveLen(2))
		})

		It("should return empty slice when no trades since time", func() {
			ctx := context.Background()

			since := time.Now()

			result := trades.GetTradesSince(ctx, since)

			Expect(result).To(BeEmpty())
		})
	})

	Describe("GetTradeByID", func() {
		It("should return trade for existing ID", func() {
			ctx := context.Background()

			tradeID := "trade-123"

			// Add trade to store
			store.AddTrade(connector.Trade{
				ID:    tradeID,
				Price: numerical.NewFromFloat(50000),
			})

			result := trades.GetTradeByID(ctx, tradeID)

			Expect(result).NotTo(BeNil())
			Expect(result.ID).To(Equal(tradeID))
		})

		It("should return nil for unknown ID", func() {
			ctx := context.Background()

			result := trades.GetTradeByID(ctx, "unknown")

			Expect(result).To(BeNil())
		})
	})

	Describe("GetTradeCount", func() {
		It("should return count from underlying store", func() {
			ctx := context.Background()

			// Add trades to store
			store.AddTrade(connector.Trade{ID: "trade-1"})
			store.AddTrade(connector.Trade{ID: "trade-2"})
			store.AddTrade(connector.Trade{ID: "trade-3"})

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(3))
		})

		It("should return zero when no trades exist", func() {
			ctx := context.Background()

			result := trades.GetTradeCount(ctx)

			Expect(result).To(Equal(0))
		})
	})

	Describe("GetTotalVolume", func() {
		It("should return volume from underlying store", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("BTC")

			// Add trades with quantities
			store.AddTrade(connector.Trade{
				Symbol:   "BTC",
				Quantity: numerical.NewFromFloat(5.5),
				ID:       "trade-1",
			})
			store.AddTrade(connector.Trade{
				Symbol:   "BTC",
				Quantity: numerical.NewFromFloat(3.0),
				ID:       "trade-2",
			})
			store.AddTrade(connector.Trade{
				Symbol:   "BTC",
				Quantity: numerical.NewFromFloat(2.0),
				ID:       "trade-3",
			})

			result := trades.GetTotalVolume(ctx, asset)

			Expect(result.Equal(numerical.NewFromFloat(10.5))).To(BeTrue())
		})

		It("should return zero volume when no trades for asset", func() {
			ctx := context.Background()

			asset := portfolio.NewAsset("UNKNOWN")

			result := trades.GetTotalVolume(ctx, asset)

			Expect(result.IsZero()).To(BeTrue())
		})
	})
})
