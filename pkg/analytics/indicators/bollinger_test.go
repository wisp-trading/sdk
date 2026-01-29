package indicators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

var _ = Describe("BollingerBands", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data for period", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.Upper.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept exactly period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.GreaterThan(result.Middle)).To(BeTrue())
			Expect(result.Middle.GreaterThan(result.Lower)).To(BeTrue())
		})

		It("should handle negative stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, -2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.LessThan(result.Middle)).To(BeTrue())
		})
	})

	Describe("Band calculation", func() {
		It("should calculate bands for simple data", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.GreaterThan(result.Middle)).To(BeTrue())
			Expect(result.Middle.GreaterThan(result.Lower)).To(BeTrue())
		})

		It("should have middle band equal to SMA", func() {
			prices := makeDecimals(120, 130, 140)
			period := 3

			result, err := indicators.BollingerBands(prices, period, 2.0)

			Expect(err).NotTo(HaveOccurred())
			expectedMiddle := numerical.NewFromFloat(130.0)
			Expect(result.Middle.Equal(expectedMiddle)).To(BeTrue())
		})

		It("should handle different stdDev multipliers", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result1, _ := indicators.BollingerBands(prices, 3, 1.0)
			result2, _ := indicators.BollingerBands(prices, 3, 2.0)
			result3, _ := indicators.BollingerBands(prices, 3, 3.0)

			Expect(result1.Middle.Equal(result2.Middle)).To(BeTrue())
			Expect(result2.Middle.Equal(result3.Middle)).To(BeTrue())

			Expect(result1.Upper.LessThan(result2.Upper)).To(BeTrue())
			Expect(result2.Upper.LessThan(result3.Upper)).To(BeTrue())

			Expect(result1.Lower.GreaterThan(result2.Lower)).To(BeTrue())
			Expect(result2.Lower.GreaterThan(result3.Lower)).To(BeTrue())
		})

		It("should handle different periods correctly", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145)

			result5, err5 := indicators.BollingerBands(prices, 5, 2.0)
			result10, err10 := indicators.BollingerBands(prices, 10, 2.0)

			Expect(err5).NotTo(HaveOccurred())
			Expect(err10).NotTo(HaveOccurred())
			Expect(result5.Middle.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(result10.Middle.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Band symmetry", func() {
		It("should have symmetric bands around middle", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			upperDist := result.Upper.Sub(result.Middle)
			lowerDist := result.Middle.Sub(result.Lower)
			Expect(upperDist.Equal(lowerDist)).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices (zero volatility)", func() {
			prices := makeDecimals(100, 100, 100, 100, 100)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.Equal(result.Middle)).To(BeTrue())
			Expect(result.Middle.Equal(result.Lower)).To(BeTrue())
			Expect(result.Middle.Equal(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle single period", func() {
			prices := makeDecimals(100, 105, 110)

			result, err := indicators.BollingerBands(prices, 1, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Middle.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			bandwidth := result.Upper.Sub(result.Lower)
			Expect(bandwidth.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle fractional prices", func() {
			prices := makeDecimals(100.5, 101.75, 102.25, 103.5, 104.0)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.GreaterThan(result.Middle)).To(BeTrue())
			Expect(result.Middle.GreaterThan(result.Lower)).To(BeTrue())
		})

		It("should handle very small stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 0.1)

			Expect(err).NotTo(HaveOccurred())
			upperDist := result.Upper.Sub(result.Middle)
			lowerDist := result.Middle.Sub(result.Lower)

			upperDistFloat, _ := upperDist.Float64()
			lowerDistFloat, _ := lowerDist.Float64()

			Expect(upperDistFloat).To(BeNumerically("<", 5.0))
			Expect(lowerDistFloat).To(BeNumerically("<", 5.0))
		})

		It("should handle very large stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 10.0)

			Expect(err).NotTo(HaveOccurred())
			bandwidth := result.Upper.Sub(result.Lower)
			bandwidthFloat, _ := bandwidth.Float64()
			Expect(bandwidthFloat).To(BeNumerically(">", 50.0))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate bands for typical market data", func() {
			prices := makeDecimals(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
				50850, 51000, 50950, 51100, 51050,
			)

			result, err := indicators.BollingerBands(prices, 20, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.GreaterThan(result.Middle)).To(BeTrue())
			Expect(result.Middle.GreaterThan(result.Lower)).To(BeTrue())
			Expect(result.Middle.GreaterThan(numerical.NewFromFloat(50000.0))).To(BeTrue())
		})

		It("should handle trending market", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				prices[i] = 100.0 + float64(i)*2
			}

			result, err := indicators.BollingerBands(prices, 10, 2.0)

			Expect(err).NotTo(HaveOccurred())
			middleVal, _ := result.Middle.Float64()
			Expect(middleVal).To(BeNumerically(">", 100.0))
		})

		It("should handle ranging market", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				prices[i] = 105.0 + float64(i%2)*5 - 2.5
			}

			result, err := indicators.BollingerBands(prices, 10, 2.0)

			Expect(err).NotTo(HaveOccurred())
			middleFloat, _ := result.Middle.Float64()
			Expect(middleFloat).To(BeNumerically("~", 105.0, 10.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 20-period, 2-stddev parameters", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.BollingerBands(prices, 20, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Upper.GreaterThan(result.Middle)).To(BeTrue())
			Expect(result.Middle.GreaterThan(result.Lower)).To(BeTrue())
		})
	})
})
