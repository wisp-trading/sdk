package indicators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("MACD", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data", func() {
			prices := makeDecimals(100, 105, 110)

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.MACD.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept minimum required data", func() {
			prices := make([]float64, 35)
			for i := 0; i < 35; i++ {
				prices[i] = 100.0 + float64(i)
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
		})
	})

	Describe("MACD calculation", func() {
		It("should calculate MACD, Signal, and Histogram", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i)
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
			Expect(result.Signal).NotTo(Equal(numerical.Zero()))
			Expect(result.Histogram).NotTo(Equal(numerical.Zero()))
		})

		It("should have histogram equal to MACD minus Signal", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i)*2
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			expectedHistogram := result.MACD.Sub(result.Signal)

			histFloat, _ := result.Histogram.Float64()
			expectedFloat, _ := expectedHistogram.Float64()

			Expect(histFloat).To(BeNumerically("~", expectedFloat, 0.01))
		})

		It("should have positive MACD in uptrend", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i)*5
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should have negative MACD in downtrend", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 200.0 - float64(i)*5
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.LessThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.Equal(numerical.Zero())).To(BeTrue())
			Expect(result.Signal.Equal(numerical.Zero())).To(BeTrue())
			Expect(result.Histogram.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should handle large price swings", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				if i%2 == 0 {
					prices[i] = 100.0
				} else {
					prices[i] = 200.0
				}
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
		})

		It("should handle fractional prices", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.5 + float64(i)*0.25
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
		})

		It("should handle short periods", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				prices[i] = 100.0 + float64(i)
			}

			result, err := indicators.MACD(prices, 5, 10, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
		})
	})

	Describe("Signal crossovers", func() {
		It("should detect bullish crossover potential", func() {
			prices := make([]float64, 50)
			for i := 0; i < 40; i++ {
				prices[i] = 100.0
			}
			for i := 40; i < 50; i++ {
				prices[i] = 100.0 + float64(i-39)*5
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.GreaterThan(result.Signal)).To(BeTrue())
		})

		It("should detect bearish crossover potential", func() {
			prices := make([]float64, 50)
			for i := 0; i < 40; i++ {
				prices[i] = 150.0
			}
			for i := 40; i < 50; i++ {
				prices[i] = 150.0 - float64(i-39)*5
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.LessThan(result.Signal)).To(BeTrue())
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate MACD for typical market data", func() {
			prices := makeDecimals(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
				50850, 51000, 50950, 51100, 51050,
				51200, 51150, 51300, 51250, 51400,
				51350, 51500, 51450, 51600, 51550,
				51700, 51650, 51800, 51750, 51900,
			)

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
			Expect(result.Signal).NotTo(Equal(numerical.Zero()))
		})

		It("should handle uptrend", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i)*2
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle downtrend", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 200.0 - float64(i)*2
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD.LessThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle ranging market", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				price := 105.0 + float64(i%2)*5 - 2.5
				prices[i] = price
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			macdFloat, _ := result.MACD.Float64()
			Expect(macdFloat).To(BeNumerically("~", 0.0, 5.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 12,26,9 parameters", func() {
			prices := make([]float64, 100)
			for i := 0; i < 100; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.MACD(prices, 12, 26, 9)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.MACD).NotTo(Equal(numerical.Zero()))
			Expect(result.Signal).NotTo(Equal(numerical.Zero()))
			Expect(result.Histogram).NotTo(Equal(numerical.Zero()))
		})
	})
})
