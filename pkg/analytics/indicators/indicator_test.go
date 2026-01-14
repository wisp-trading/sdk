package indicators_test

import (
	"context"
	"time"

	sdkTesting "github.com/backtesting-org/kronos-sdk/pkg/testing"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	marketTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("Indicators with Market Registry", func() {
	var (
		app           *fxtest.App
		indicatorsSvc analytics.Indicators
		registry      marketTypes.MarketRegistry
		ctx           context.Context
		btc           portfolio.Asset
		exchangeName  connector.ExchangeName
	)

	BeforeEach(func() {
		btc = portfolio.NewAsset("BTC")
		exchangeName = "binance"
		ctx = context.Background()

		app = fxtest.New(GinkgoT(),
			sdkTesting.Module,
			fx.Populate(
				&indicatorsSvc,
				&registry,
			),
			fx.NopLogger,
		)

		app.RequireStart()

		// Get spot store from registry
		spotStore := registry.Get(marketTypes.MarketTypeSpot)
		Expect(spotStore).ToNot(BeNil())

		// Populate test data - add klines to the store
		populateTestKlines(spotStore, btc, exchangeName, 100)
	})

	AfterEach(func() {
		app.RequireStop()
	})

	Context("RSI indicator", func() {
		It("should calculate RSI using market registry", func() {
			result, err := indicatorsSvc.RSI(ctx, btc, 14, analytics.IndicatorOptions{
				Exchange: exchangeName,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("SMA indicator", func() {
		It("should calculate SMA using market registry", func() {
			result, err := indicatorsSvc.SMA(ctx, btc, 20, analytics.IndicatorOptions{
				Exchange: exchangeName,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("EMA indicator", func() {
		It("should calculate EMA using market registry", func() {
			result, err := indicatorsSvc.EMA(ctx, btc, 50, analytics.IndicatorOptions{
				Exchange: exchangeName,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("Multiple indicator calls", func() {
		It("should handle multiple indicator calculations", func() {
			rsi, err := indicatorsSvc.RSI(ctx, btc, 14, analytics.IndicatorOptions{Exchange: exchangeName})
			Expect(err).ToNot(HaveOccurred())

			sma, err := indicatorsSvc.SMA(ctx, btc, 20, analytics.IndicatorOptions{Exchange: exchangeName})
			Expect(err).ToNot(HaveOccurred())

			ema, err := indicatorsSvc.EMA(ctx, btc, 50, analytics.IndicatorOptions{Exchange: exchangeName})
			Expect(err).ToNot(HaveOccurred())

			Expect(rsi.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(sma.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(ema.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})
})

// populateTestKlines adds test klines to the store
func populateTestKlines(store marketTypes.MarketStore, asset portfolio.Asset, exchange connector.ExchangeName, count int) {
	now := time.Now()
	// Use the default interval that matches what indicators use
	interval := analytics.DefaultInterval

	for i := 0; i < count; i++ {
		kline := connector.Kline{
			Interval:  interval,
			OpenTime:  now.Add(time.Duration(-count+i) * time.Hour),
			Open:      float64(100 + i),
			High:      float64(105 + i),
			Low:       float64(95 + i),
			Close:     float64(100 + i),
			Volume:    1000.0,
			CloseTime: now.Add(time.Duration(-count+i+1) * time.Hour),
		}
		// Add klines to store
		store.UpdateKline(asset, exchange, kline)
	}

	// Add a price so getDefaultExchange works
	store.UpdateAssetPrice(asset, exchange, connector.Price{
		Price:     numerical.NewFromInt(100),
		Timestamp: now,
	})
}
