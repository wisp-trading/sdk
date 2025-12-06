package indicators_test

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SMA", func() {
	Describe("Input validation", func() {
		It("should return error when insufficient data", func() {
			prices := makeDecimals(100, 105)

			result, err := indicators.SMA(prices, 14)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept exactly period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should accept more than period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("SMA calculation", func() {
		It("should calculate SMA for simple data", func() {
			prices := makeDecimals(100, 110, 120, 130, 140)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			expected := numerical.NewFromFloat(120.0) // (100+110+120+130+140)/5
			Expect(result.Equal(expected)).To(BeTrue())
		})

		It("should calculate SMA using most recent prices", func() {
			prices := makeDecimals(100, 110, 120, 130, 140, 150)

			result, err := indicators.SMA(prices, 3)

			Expect(err).NotTo(HaveOccurred())
			expected := numerical.NewFromFloat(140.0) // (130+140+150)/3
			Expect(result.Equal(expected)).To(BeTrue())
		})

		It("should handle different periods correctly", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145)

			sma3, _ := indicators.SMA(prices, 3)
			sma5, _ := indicators.SMA(prices, 5)
			sma10, _ := indicators.SMA(prices, 10)

			Expect(sma3.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(sma5.GreaterThan(numerical.Zero())).To(BeTrue())
			Expect(sma10.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should calculate correct mathematical average", func() {
			prices := makeDecimals(10, 20, 30, 40, 50)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			expected := numerical.NewFromFloat(30.0)
			Expect(result.Equal(expected)).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices", func() {
			prices := makeDecimals(100, 100, 100, 100, 100)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Equal(numerical.NewFromFloat(100.0))).To(BeTrue())
		})

		It("should handle single period", func() {
			prices := makeDecimals(100, 105, 110)

			result, err := indicators.SMA(prices, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Equal(numerical.NewFromFloat(110.0))).To(BeTrue())
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			expected := numerical.NewFromFloat(135.0) // (100+200+50+300+25)/5
			Expect(result.Equal(expected)).To(BeTrue())
		})

		It("should handle fractional prices", func() {
			prices := makeDecimalsFloat(100.5, 101.75, 102.25, 103.5, 104.0)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(102.0))).To(BeTrue())
		})

		It("should handle very small period", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.SMA(prices, 2)

			Expect(err).NotTo(HaveOccurred())
			expected := numerical.NewFromFloat(117.5) // (115+120)/2
			Expect(result.Equal(expected)).To(BeTrue())
		})

		It("should handle very large period", func() {
			prices := make([]numerical.Decimal, 200)
			for i := 0; i < 200; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i))
			}

			result, err := indicators.SMA(prices, 100)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(100.0))).To(BeTrue())
		})
	})

	Describe("Trend identification", func() {
		It("should reflect uptrend", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i)*5)
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			// Latest 10 prices average should be higher than starting price
			Expect(result.GreaterThan(numerical.NewFromFloat(150.0))).To(BeTrue())
		})

		It("should reflect downtrend", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(200.0 - float64(i)*5)
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			// Latest 10 prices average should be lower than starting price
			Expect(result.LessThan(numerical.NewFromFloat(150.0))).To(BeTrue())
		})

		It("should smooth out noise", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 20; i++ {
				base := 105.0
				if i%2 == 0 {
					base += 5.0
				} else {
					base -= 5.0
				}
				prices[i] = numerical.NewFromFloat(base)
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 105.0, 5.0))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate SMA for typical market data", func() {
			prices := makeDecimalsFloat(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
				50850, 51000, 50950, 51100, 51050,
			)

			result, err := indicators.SMA(prices, 20)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromFloat(50000.0))).To(BeTrue())
			Expect(result.LessThan(numerical.NewFromFloat(52000.0))).To(BeTrue())
		})

		It("should handle volatile market", func() {
			prices := make([]numerical.Decimal, 30)
			for i := 0; i < 30; i++ {
				base := 100.0
				if i%3 == 0 {
					base += 10.0
				} else if i%3 == 1 {
					base -= 10.0
				}
				prices[i] = numerical.NewFromFloat(base)
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 100.0, 15.0))
		})

		It("should handle ranging market", func() {
			prices := make([]numerical.Decimal, 30)
			for i := 0; i < 30; i++ {
				price := 105.0 + float64(i%2)*5 - 2.5
				prices[i] = numerical.NewFromFloat(price)
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 105.0, 5.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 20-period SMA", func() {
			prices := make([]numerical.Decimal, 50)
			for i := 0; i < 50; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i%10))
			}

			result, err := indicators.SMA(prices, 20)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should work with standard 50-period SMA", func() {
			prices := make([]numerical.Decimal, 100)
			for i := 0; i < 100; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i%10))
			}

			result, err := indicators.SMA(prices, 50)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should work with standard 200-period SMA", func() {
			prices := make([]numerical.Decimal, 400)
			for i := 0; i < 400; i++ {
				prices[i] = numerical.NewFromFloat(100.0 + float64(i)*0.1)
			}

			result, err := indicators.SMA(prices, 200)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Comparison with EMA", func() {
		It("should be less responsive than EMA to recent price changes", func() {
			prices := make([]numerical.Decimal, 20)
			for i := 0; i < 15; i++ {
				prices[i] = numerical.NewFromFloat(100.0)
			}
			for i := 15; i < 20; i++ {
				prices[i] = numerical.NewFromFloat(150.0)
			}

			sma, _ := indicators.SMA(prices, 10)
			ema, _ := indicators.EMA(prices, 10)

			smaFloat, _ := sma.Float64()
			emaFloat, _ := ema.Float64()

			Expect(emaFloat).To(BeNumerically(">", smaFloat))
		})
	})
})
