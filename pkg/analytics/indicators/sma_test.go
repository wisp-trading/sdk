package indicators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/sdk/pkg/analytics/indicators"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
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
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 0.0))
		})

		It("should accept more than period data points", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 0.0))
		})
	})

	Describe("SMA calculation", func() {
		It("should calculate SMA for simple data", func() {
			prices := makeDecimals(100, 110, 120, 130, 140)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 120.0, 0.01)) // (100+110+120+130+140)/5
		})

		It("should calculate SMA using most recent prices", func() {
			prices := makeDecimals(100, 110, 120, 130, 140, 150)

			result, err := indicators.SMA(prices, 3)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 140.0, 0.01)) // (130+140+150)/3
		})

		It("should handle different periods correctly", func() {
			prices := makeDecimals(100, 105, 110, 115, 120, 125, 130, 135, 140, 145)

			sma3, _ := indicators.SMA(prices, 3)
			sma5, _ := indicators.SMA(prices, 5)
			sma10, _ := indicators.SMA(prices, 10)

			val3, _ := sma3.Float64()
			val5, _ := sma5.Float64()
			val10, _ := sma10.Float64()
			Expect(val3).To(BeNumerically(">", 0.0))
			Expect(val5).To(BeNumerically(">", 0.0))
			Expect(val10).To(BeNumerically(">", 0.0))
		})

		It("should calculate correct mathematical average", func() {
			prices := makeDecimals(10, 20, 30, 40, 50)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 30.0, 0.01))
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices", func() {
			prices := makeDecimals(100, 100, 100, 100, 100)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 100.0, 0.01))
		})

		It("should handle single period", func() {
			prices := makeDecimals(100, 105, 110)

			result, err := indicators.SMA(prices, 1)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 110.0, 0.01))
		})

		It("should handle large price swings", func() {
			prices := makeDecimals(100, 200, 50, 300, 25)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 135.0, 0.01)) // (100+200+50+300+25)/5
		})

		It("should handle fractional prices", func() {
			prices := makeDecimals(100.5, 101.75, 102.25, 103.5, 104.0)

			result, err := indicators.SMA(prices, 5)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 102.0))
		})

		It("should handle very small period", func() {
			prices := makeDecimals(100, 105, 110, 115, 120)

			result, err := indicators.SMA(prices, 2)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("~", 117.5, 0.01)) // (115+120)/2
		})

		It("should handle very large period", func() {
			prices := make([]float64, 200)
			for i := 0; i < 200; i++ {
				prices[i] = 100.0 + float64(i)
			}

			result, err := indicators.SMA(prices, 100)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 100.0))
		})
	})

	Describe("Trend identification", func() {
		It("should reflect uptrend", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				prices[i] = 100.0 + float64(i)*5
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 150.0))
		})

		It("should reflect downtrend", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				prices[i] = 200.0 - float64(i)*5
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically("<", 150.0))
		})

		It("should smooth out noise", func() {
			prices := make([]float64, 20)
			for i := 0; i < 20; i++ {
				base := 105.0
				if i%2 == 0 {
					base += 5.0
				} else {
					base -= 5.0
				}
				prices[i] = base
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 105.0, 5.0))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate SMA for typical market data", func() {
			prices := makeDecimals(
				50100, 50250, 50150, 50300, 50400,
				50350, 50500, 50450, 50600, 50550,
				50700, 50650, 50800, 50750, 50900,
				50850, 51000, 50950, 51100, 51050,
			)

			result, err := indicators.SMA(prices, 20)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 50000.0))
			Expect(smaVal).To(BeNumerically("<", 52000.0))
		})

		It("should handle volatile market", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				base := 100.0
				if i%3 == 0 {
					base += 10.0
				} else if i%3 == 1 {
					base -= 10.0
				}
				prices[i] = base
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 100.0, 15.0))
		})

		It("should handle ranging market", func() {
			prices := make([]float64, 30)
			for i := 0; i < 30; i++ {
				price := 105.0 + float64(i%2)*5 - 2.5
				prices[i] = price
			}

			result, err := indicators.SMA(prices, 10)

			Expect(err).NotTo(HaveOccurred())
			resultFloat, _ := result.Float64()
			Expect(resultFloat).To(BeNumerically("~", 105.0, 5.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 20-period SMA", func() {
			prices := make([]float64, 50)
			for i := 0; i < 50; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.SMA(prices, 20)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 0.0))
		})

		It("should work with standard 50-period SMA", func() {
			prices := make([]float64, 100)
			for i := 0; i < 100; i++ {
				prices[i] = 100.0 + float64(i%10)
			}

			result, err := indicators.SMA(prices, 50)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 0.0))
		})

		It("should work with standard 200-period SMA", func() {
			prices := make([]float64, 400)
			for i := 0; i < 400; i++ {
				prices[i] = 100.0 + float64(i)*0.1
			}

			result, err := indicators.SMA(prices, 200)

			Expect(err).NotTo(HaveOccurred())
			smaVal, _ := result.Float64()
			Expect(smaVal).To(BeNumerically(">", 0.0))
		})
	})

	Describe("Comparison with EMA", func() {
		It("should be less responsive than EMA to recent price changes", func() {
			prices := make([]float64, 20)
			for i := 0; i < 15; i++ {
				prices[i] = 100.0
			}
			for i := 15; i < 20; i++ {
				prices[i] = 150.0
			}

			sma, _ := indicators.SMA(prices, 10)
			ema, _ := indicators.EMA(prices, 10)

			smaFloat, _ := sma.Float64()
			emaFloat, _ := ema.Float64()

			Expect(emaFloat).To(BeNumerically(">", smaFloat))
		})
	})
})
