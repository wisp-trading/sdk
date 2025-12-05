package indicators_test

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BollingerBands", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data for period", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result).To(BeNil())
		})

		It("should accept exactly period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(1))
		})

		It("should handle negative stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, -2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("Band calculation", func() {
		It("should calculate bands for simple data", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
			Expect(len(result)).To(Equal(3))

			// Each result should have upper > middle > lower
			for _, bb := range result {
				Expect(bb.Upper.GreaterThan(bb.Middle)).To(BeTrue())
				Expect(bb.Middle.GreaterThan(bb.Lower)).To(BeTrue())
			}
		})

		It("should have middle band equal to SMA", func() {
			prices := makeDecimals(100, 110, 120, 130, 140)
			period := 3

			result, err := indicators.BollingerBands(prices, period, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// First middle band should be SMA of first 3 prices: (100+110+120)/3 = 110
			expectedMiddle := numerical.NewFromFloat(110.0)
			Expect(result[0].Middle.Equal(expectedMiddle)).To(BeTrue())
		})

		It("should handle different stdDev multipliers", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result1, _ := indicators.BollingerBands(prices, 3, 1.0)
			result2, _ := indicators.BollingerBands(prices, 3, 2.0)
			result3, _ := indicators.BollingerBands(prices, 3, 3.0)

			Expect(len(result1)).To(Equal(len(result2)))
			Expect(len(result2)).To(Equal(len(result3)))

			// Middle bands should be the same
			Expect(result1[0].Middle.Equal(result2[0].Middle)).To(BeTrue())
			Expect(result2[0].Middle.Equal(result3[0].Middle)).To(BeTrue())

			// Upper band should increase with larger multiplier
			Expect(result1[0].Upper.LessThan(result2[0].Upper)).To(BeTrue())
			Expect(result2[0].Upper.LessThan(result3[0].Upper)).To(BeTrue())

			// Lower band should decrease with larger multiplier
			Expect(result1[0].Lower.GreaterThan(result2[0].Lower)).To(BeTrue())
			Expect(result2[0].Lower.GreaterThan(result3[0].Lower)).To(BeTrue())
		})

		It("should handle different periods correctly", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145)

			result5, _ := indicators.BollingerBands(prices, 5, 2.0)
			result10, _ := indicators.BollingerBands(prices, 10, 2.0)

			// Longer period = fewer results
			Expect(len(result5)).To(BeNumerically(">", len(result10)))
			Expect(len(result5)).To(Equal(6))
			Expect(len(result10)).To(Equal(1))
		})
	})

	Describe("Band symmetry", func() {
		It("should have symmetric bands around middle", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())

			// Distance from upper to middle should equal distance from middle to lower
			for _, bb := range result {
				upperDist := bb.Upper.Sub(bb.Middle)
				lowerDist := bb.Middle.Sub(bb.Lower)
				Expect(upperDist.Equal(lowerDist)).To(BeTrue())
			}
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices (zero volatility)", func() {
			prices := makeDecimals(100, 100, 100, 100, 100)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// With zero volatility, all bands should be equal to price
			for _, bb := range result {
				Expect(bb.Upper.Equal(bb.Middle)).To(BeTrue())
				Expect(bb.Middle.Equal(bb.Lower)).To(BeTrue())
				Expect(bb.Middle.Equal(numerical.NewFromFloat(100.0))).To(BeTrue())
			}
		})

		It("should handle single period", func() {
			prices := makeDecimals(100, 105, 110)

			result, err := indicators.BollingerBands(prices, 1, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
			Expect(len(result)).To(Equal(3))
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// Bands should be wide due to high volatility
			for _, bb := range result {
				bandwidth := bb.Upper.Sub(bb.Lower)
				Expect(bandwidth.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
			}
		})

		It("should handle fractional prices", func() {
			prices := makeDecimalsFloat(100.5, 101.75, 102.25, 103.5, 104.0)

			result, err := indicators.BollingerBands(prices, 3, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			for _, bb := range result {
				Expect(bb.Upper.GreaterThan(bb.Middle)).To(BeTrue())
				Expect(bb.Middle.GreaterThan(bb.Lower)).To(BeTrue())
			}
		})

		It("should handle very small stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 0.1)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// Bands should be very close to middle
			for _, bb := range result {
				upperDist := bb.Upper.Sub(bb.Middle)
				lowerDist := bb.Middle.Sub(bb.Lower)

				upperDistFloat, _ := upperDist.Float64()
				lowerDistFloat, _ := lowerDist.Float64()

				Expect(upperDistFloat).To(BeNumerically("<", 5.0))
				Expect(lowerDistFloat).To(BeNumerically("<", 5.0))
			}
		})

		It("should handle very large stdDev multiplier", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.BollingerBands(prices, 3, 10.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// Bands should be very wide
			for _, bb := range result {
				bandwidth := bb.Upper.Sub(bb.Lower)
				bandwidthFloat, _ := bandwidth.Float64()
				Expect(bandwidthFloat).To(BeNumerically(">", 50.0))
			}
		})
	})

	Describe("Result length", func() {
		It("should return correct number of values", func() {
			dataLength := 100
			period := 20

			prices := make([]numerical.Decimal, dataLength)
			for i := 0; i < dataLength; i++ {
				prices[i] = numerical.NewFromInt(int64(100 + i))
			}

			result, err := indicators.BollingerBands(prices, period, 2.0)

			Expect(err).NotTo(HaveOccurred())
			expectedLength := dataLength - period + 1
			Expect(len(result)).To(Equal(expectedLength))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate bands for typical market data", func() {
			// Simulating realistic price data
			prices := makeDecimalsFloat(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
				50850, 51000, 50950, 51100, 51050,
			)

			result, err := indicators.BollingerBands(prices, 20, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
			Expect(len(result)).To(Equal(1))

			// All bands should be reasonable
			bb := result[0]
			Expect(bb.Upper.GreaterThan(bb.Middle)).To(BeTrue())
			Expect(bb.Middle.GreaterThan(bb.Lower)).To(BeTrue())
			Expect(bb.Middle.GreaterThan(numerical.NewFromFloat(50000.0))).To(BeTrue())
		})

		It("should handle trending market", func() {
			prices := []numerical.Decimal{}

			// Simulate uptrend
			for i := 0; i < 30; i++ {
				price := 100.0 + float64(i)*2
				prices = append(prices, numerical.NewFromFloat(price))
			}

			result, err := indicators.BollingerBands(prices, 10, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// Middle band should trend upward
			firstMiddle := result[0].Middle
			lastMiddle := result[len(result)-1].Middle
			Expect(lastMiddle.GreaterThan(firstMiddle)).To(BeTrue())
		})

		It("should handle ranging market", func() {
			prices := []numerical.Decimal{}

			// Simulate ranging market between 100 and 110
			for i := 0; i < 30; i++ {
				price := 105.0 + float64(i%2)*5 - 2.5
				prices = append(prices, numerical.NewFromFloat(price))
			}

			result, err := indicators.BollingerBands(prices, 10, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			// Bands should be relatively stable
			for _, bb := range result {
				middleFloat, _ := bb.Middle.Float64()
				Expect(middleFloat).To(BeNumerically("~", 105.0, 10.0))
			}
		})

		It("should squeeze during low volatility", func() {
			prices := []numerical.Decimal{}

			// Start with volatility
			for i := 0; i < 10; i++ {
				price := 100.0 + float64(i%2)*10
				prices = append(prices, numerical.NewFromFloat(price))
			}
			// Then reduce volatility
			for i := 0; i < 10; i++ {
				price := 100.0 + float64(i%2)*2
				prices = append(prices, numerical.NewFromFloat(price))
			}

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeNumerically(">", 5))

			// Later bands should be narrower
			earlyBandwidth := result[0].Upper.Sub(result[0].Lower)
			lateBandwidth := result[len(result)-1].Upper.Sub(result[len(result)-1].Lower)
			Expect(lateBandwidth.LessThan(earlyBandwidth)).To(BeTrue())
		})

		It("should expand during high volatility", func() {
			prices := []numerical.Decimal{}

			// Start with low volatility
			for i := 0; i < 10; i++ {
				price := 100.0 + float64(i%2)*2
				prices = append(prices, numerical.NewFromFloat(price))
			}
			// Then increase volatility
			for i := 0; i < 10; i++ {
				price := 100.0 + float64(i%2)*20
				prices = append(prices, numerical.NewFromFloat(price))
			}

			result, err := indicators.BollingerBands(prices, 5, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeNumerically(">", 5))

			// Later bands should be wider
			earlyBandwidth := result[0].Upper.Sub(result[0].Lower)
			lateBandwidth := result[len(result)-1].Upper.Sub(result[len(result)-1].Lower)
			Expect(lateBandwidth.GreaterThan(earlyBandwidth)).To(BeTrue())
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 20-period, 2-stddev parameters", func() {
			prices := make([]numerical.Decimal, 50)
			for i := 0; i < 50; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i%10))
			}

			result, err := indicators.BollingerBands(prices, 20, 2.0)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
			Expect(len(result)).To(Equal(31))
		})
	})
})
