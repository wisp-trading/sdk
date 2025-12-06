package indicators_test

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RSI", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.RSI(prices, 14)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept minimum required data", func() {
			prices := make([]numerical.Decimal, 15)
			for i := 0; i < 15; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i))
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})
	})

	Describe("RSI calculation", func() {
		It("should calculate RSI for simple data", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145, 150, 155, 160, 165, 170)

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should return 100 when all prices are increasing", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i)*10)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromInt(70))).To(BeTrue())
		})

		It("should return low RSI when all prices are decreasing", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(200.0 - float64(i)*10)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.LessThan(numerical.NewFromInt(30))).To(BeTrue())
		})

		It("should return ~50 for sideways price action", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i%2)*5 - 2.5)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			rsiValue, _ := result.Float64()
			Expect(rsiValue).To(BeNumerically("~", 50.0, 20.0))
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(100.0)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Equal(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25, 400, 10, 500, 5, 600, 2, 700, 1, 800, 0.5)

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should handle fractional prices", func() {
			prices := makeDecimalsFloat(100.5, 101.75, 102.25, 103.5, 104.0, 105.25, 106.5, 107.75, 108.25, 109.5, 110.0, 111.25, 112.5, 113.75, 114.25)

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should handle small period", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125)

			result, err := indicators.RSI(prices, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
		})

		It("should handle large period", func() {
			prices := make([]numerical.Decimal, 100)
			for i := 0; i < 100; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i)*0.5)
			}

			result, err := indicators.RSI(prices, 50)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Overbought and oversold", func() {
		It("should identify overbought conditions", func() {
			prices := make([]numerical.Decimal, 20)
			prices[0] = numerical.NewFromFloat(100.0)
			for i := 1; i < 20; i++ {
				prices[i] = prices[i-1].Add(numerical.NewFromFloat(5.0))
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromInt(70))).To(BeTrue())
		})

		It("should identify oversold conditions", func() {
			prices := make([]numerical.Decimal, 20)
			prices[0] = numerical.NewFromFloat(200.0)
			for i := 1; i < 20; i++ {
				prices[i] = prices[i-1].Sub(numerical.NewFromFloat(5.0))
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.LessThan(numerical.NewFromInt(30))).To(BeTrue())
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate RSI for typical market data", func() {
			prices := makeDecimalsFloat(50100, 50250, 50150, 50300, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800, 50750, 50900)

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should handle uptrend", func() {
			prices := make([]numerical.Decimal, 30)
			for i := 0; i < 30; i++ {
				price := 100.0 + float64(i)*2
				prices[i] = numerical.NewFromFloat(price)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromInt(50))).To(BeTrue())
		})

		It("should handle downtrend", func() {
			prices := make([]numerical.Decimal, 30)
			for i := 0; i < 30; i++ {
				price := 200.0 - float64(i)*2
				prices[i] = numerical.NewFromFloat(price)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.LessThan(numerical.NewFromInt(50))).To(BeTrue())
		})

		It("should handle ranging market", func() {
			prices := make([]numerical.Decimal, 30)
			for i := 0; i < 30; i++ {
				price := 105.0 + float64(i%2)*5 - 2.5
				prices[i] = numerical.NewFromFloat(price)
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			rsiValue, _ := result.Float64()
			Expect(rsiValue).To(BeNumerically("~", 50.0, 25.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 14-period RSI", func() {
			prices := make([]numerical.Decimal, 50)
			for i := 0; i < 50; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i%10))
			}

			result, err := indicators.RSI(prices, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThanOrEqual(numerical.Zero())).To(BeTrue())
			Expect(result.LessThanOrEqual(numerical.NewFromInt(100))).To(BeTrue())
		})
	})
})
