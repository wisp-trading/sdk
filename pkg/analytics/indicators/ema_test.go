package indicators_test

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EMA", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data for period", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.EMA(prices, 5)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept exactly period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.EMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("EMA calculation", func() {
		It("should calculate EMA for simple data", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.EMA(prices, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should be more responsive than SMA", func() {
			prices := makeDecimals(100, 100, 100, 100, 150)

			ema, _ := indicators.EMA(prices, 3)
			sma, _ := indicators.SMA(prices, 3)

			emaValue, _ := ema.Float64()
			smaValue, _ := sma.Float64()

			Expect(emaValue).To(BeNumerically(">", smaValue))
		})

		It("should handle trending data", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				prices[i] = 100.0 + float64(i)*5
			}

			result, err := indicators.EMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle different periods", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145)

			ema5, _ := indicators.EMA(prices, 5)
			ema10, _ := indicators.EMA(prices, 10)

			Expect(ema5.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(ema10.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices", func() {
			prices := makeDecimals(100, 100, 100, 100, 100)

			result, err := indicators.EMA(prices, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Equal(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle single period", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.EMA(prices, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25)

			result, err := indicators.EMA(prices, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle fractional prices", func() {
			prices := makeDecimals(100.5, 101.75, 102.25, 103.5, 104.0, 105.25, 106.5, 107.75, 108.25, 109.5, 110.0, 111.25, 112.5, 113.75, 114.25)

			result, err := indicators.EMA(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle large period", func() {
			prices := make([]float64, 100)
			for i := 0; i < 100; i++ {
				prices[i] = 100.0 + float64(i)*0.5
			}

			result, err := indicators.EMA(prices, 100)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
		})
	})

	Describe("Exponential weighting", func() {
		It("should give more weight to recent prices", func() {
			prices1 := makeDecimals(100, 100, 100, 150)
			prices2 := makeDecimals(150, 100, 100, 100)

			ema1, _ := indicators.EMA(prices1, 3)
			ema2, _ := indicators.EMA(prices2, 3)

			Expect(ema1.GreaterThan(ema2)).To(BeTrue())
		})

		It("should approach recent price faster with shorter period", func() {
			prices := makeDecimals(100, 100, 100, 100, 200)

			ema5, _ := indicators.EMA(prices, 3)
			ema10, _ := indicators.EMA(prices, 4)

			lastPrice := numerical.NewFromFloat(200.0)
			diff5 := lastPrice.Sub(ema5)
			diff10 := lastPrice.Sub(ema10)

			diff5Val, _ := diff5.Float64()
			diff10Val, _ := diff10.Float64()

			Expect(diff5Val).To(BeNumerically("<", diff10Val))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate EMA for typical market data", func() {
			prices := makeDecimals(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
			)

			result, err := indicators.EMA(prices, 12)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(50000.0))).To(BeTrue())
		})

		It("should handle uptrend", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				prices[i] = 100.0 + float64(i)*2
			}

			result, err := indicators.EMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle downtrend", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				prices[i] = 200.0 - float64(i)*2
			}

			result, err := indicators.EMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.LessThan(numerical.NewFromFloat(200.0))).To(BeTrue())
		})

		It("should handle ranging market", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				prices[i] = 105.0 + float64(i%2)*5 - 2.5
			}

			result, err := indicators.EMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 105.0, 10.0))
		})

		It("should handle volatile market", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				prices[i] = 100.0 + float64(i%2)*20 - 10.0
			}

			result, err := indicators.EMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Standard periods", func() {
		It("should work with standard 12-period EMA", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.EMA(prices, 12)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should work with standard 26-period EMA", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.EMA(prices, 26)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should work with standard 50-period EMA", func() {
			prices := make([]float64, 100)
			for i := 0; i < 100; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.EMA(prices, 50)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should work with standard 200-period EMA", func() {
			prices := make([]float64, 400)
			for i := 0; i < 400; i++ {
				prices[i] = 100.0 + float64(i)*0.1
			}

			result, err := indicators.EMA(prices, 200)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})
})
