package indicators_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	wispAnalytics "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("Indicators", func() {
	var (
		svc    wispAnalytics.Indicators
		klines []connector.Kline
	)

	BeforeEach(func() {
		svc = indicators.NewIndicators()
		klines = makeKlines(100)
	})

	Context("RSI", func() {
		It("calculates RSI from klines", func() {
			result, err := svc.RSI(klines, 14)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("SMA", func() {
		It("calculates SMA from klines", func() {
			result, err := svc.SMA(klines, 20)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("EMA", func() {
		It("calculates EMA from klines", func() {
			result, err := svc.EMA(klines, 50)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Context("multiple indicators on same klines", func() {
		It("runs all three without redundant fetches", func() {
			rsi, err := svc.RSI(klines, 14)
			Expect(err).ToNot(HaveOccurred())

			sma, err := svc.SMA(klines, 20)
			Expect(err).ToNot(HaveOccurred())

			ema, err := svc.EMA(klines, 50)
			Expect(err).ToNot(HaveOccurred())

			Expect(rsi.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(sma.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(ema.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})
})

// makeKlines produces count ascending close-price klines for testing.
func makeKlines(count int) []connector.Kline {
	now := time.Now()
	klines := make([]connector.Kline, count)
	for i := 0; i < count; i++ {
		price := float64(100 + i)
		klines[i] = connector.Kline{
			Interval:  "1h",
			OpenTime:  now.Add(time.Duration(-count+i) * time.Hour),
			CloseTime: now.Add(time.Duration(-count+i+1) * time.Hour),
			Open:      price,
			High:      price + 5,
			Low:       price - 5,
			Close:     price,
			Volume:    1000.0,
		}
	}
	return klines
}
