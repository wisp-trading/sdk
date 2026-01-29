package indicators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wisp-trading/wisp/pkg/analytics/indicators"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

var _ = Describe("Stochastic", func() {
	Describe("Input validation", func() {
		It("should return error when arrays have different lengths", func() {
			highs := makeDecimals(100, 105, 110)
			lows := makeDecimals(95, 100)
			closes := makeDecimals(98, 103, 108)

			result, err := indicators.Stochastic(highs, lows, closes, 2, 3)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("equal length"))
			Expect(result.K.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should return error when insufficient data", func() {
			highs := makeDecimals(100, 105)
			lows := makeDecimals(95, 100)
			closes := makeDecimals(98, 103)

			result, err := indicators.Stochastic(highs, lows, closes, 5, 3)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.K.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should accept minimum required data", func() {
			highs := makeDecimals(100, 105, 110, 115)
			lows := makeDecimals(95, 100, 105, 110)
			closes := makeDecimals(98, 103, 108, 113)

			result, err := indicators.Stochastic(highs, lows, closes, 2, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
		})
	})

	Describe("Stochastic calculation", func() {
		It("should calculate %K and %D for simple data", func() {
			highs := makeDecimals(100, 105, 110, 115, 120)
			lows := makeDecimals(95, 100, 105, 110, 115)
			closes := makeDecimals(98, 103, 108, 113, 118)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			dVal, _ := result.D.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
			Expect(dVal).To(BeNumerically(">=", 0.0))
			Expect(dVal).To(BeNumerically("<=", 100.0))
		})

		It("should return 100 when close equals highest high", func() {
			highs := makeDecimals(100, 105, 110, 115, 120)
			lows := makeDecimals(95, 100, 105, 110, 115)
			closes := makeDecimals(98, 103, 108, 113, 120)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.K.Equal(numerical.NewFromInt(100))).To(BeTrue())
		})

		It("should return 0 when close equals lowest low", func() {
			highs := makeDecimals(120, 115, 110, 105, 100)
			lows := makeDecimals(115, 110, 105, 100, 95)
			closes := makeDecimals(118, 113, 108, 103, 95)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.K.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should return 50 when close is at midpoint", func() {
			highs := makeDecimals(100, 100, 100, 100, 100)
			lows := makeDecimals(90, 90, 90, 90, 90)
			closes := makeDecimals(95, 95, 95, 95, 95)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.K.Equal(numerical.NewFromInt(50))).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices (no range)", func() {
			highs := makeDecimals(100, 100, 100, 100)
			lows := makeDecimals(100, 100, 100, 100)
			closes := makeDecimals(100, 100, 100, 100)

			result, err := indicators.Stochastic(highs, lows, closes, 2, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.K.Equal(numerical.NewFromInt(50))).To(BeTrue())
		})

		It("should handle large price swings", func() {
			highs := makeDecimals(100, 200, 150, 300, 250)
			lows := makeDecimals(90, 180, 140, 280, 240)
			closes := makeDecimals(95, 190, 145, 290, 245)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
		})

		It("should handle fractional prices", func() {
			highs := makeDecimals(100.5, 101.75, 102.25, 103.5)
			lows := makeDecimals(99.25, 100.5, 101.0, 102.25)
			closes := makeDecimals(100.0, 101.25, 101.75, 103.0)

			result, err := indicators.Stochastic(highs, lows, closes, 2, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
		})

		It("should handle single kPeriod", func() {
			highs := makeDecimals(100, 105, 110, 115)
			lows := makeDecimals(95, 100, 105, 110)
			closes := makeDecimals(98, 103, 108, 113)

			result, err := indicators.Stochastic(highs, lows, closes, 1, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
		})

		It("should handle single dPeriod", func() {
			highs := makeDecimals(100, 105, 110, 115)
			lows := makeDecimals(95, 100, 105, 110)
			closes := makeDecimals(98, 103, 108, 113)

			result, err := indicators.Stochastic(highs, lows, closes, 2, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.K.Equal(result.D)).To(BeTrue())
		})
	})

	Describe("Overbought and oversold", func() {
		It("should identify overbought conditions", func() {
			highs := makeDecimals(100, 105, 110, 115, 120)
			lows := makeDecimals(95, 100, 105, 110, 115)
			closes := makeDecimals(98, 103, 108, 114, 119)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">", 70.0))
		})

		It("should identify oversold conditions", func() {
			highs := makeDecimals(120, 115, 110, 105, 100)
			lows := makeDecimals(115, 110, 105, 100, 95)
			closes := makeDecimals(119, 114, 109, 104, 96)

			result, err := indicators.Stochastic(highs, lows, closes, 3, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically("<", 30.0))
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate stochastic for typical market data", func() {
			highs := makeDecimals(50100, 50250, 50150, 50300, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800, 50750, 50900, 50850, 51000)
			lows := makeDecimals(49900, 50050, 49950, 50100, 50200, 50150, 50300, 50250, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800)
			closes := makeDecimals(50000, 50150, 50050, 50200, 50300, 50250, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800, 50750, 50900)

			result, err := indicators.Stochastic(highs, lows, closes, 14, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			dVal, _ := result.D.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
			Expect(dVal).To(BeNumerically(">=", 0.0))
			Expect(dVal).To(BeNumerically("<=", 100.0))
		})

		It("should handle uptrend", func() {
			highs := []float64{}
			lows := []float64{}
			closes := []float64{}

			for i := 0; i < 30; i++ {
				base := 100.0 + float64(i)*2
				highs = append(highs, base+5)
				lows = append(lows, base-5)
				closes = append(closes, base+3)
			}

			result, err := indicators.Stochastic(highs, lows, closes, 14, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically(">", 50.0))
		})

		It("should handle downtrend", func() {
			highs := []float64{}
			lows := []float64{}
			closes := []float64{}

			for i := 0; i < 30; i++ {
				base := 200.0 - float64(i)*2
				highs = append(highs, base+5)
				lows = append(lows, base-5)
				closes = append(closes, base-3)
			}

			result, err := indicators.Stochastic(highs, lows, closes, 14, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			Expect(kVal).To(BeNumerically("<", 50.0))
		})

		It("should handle ranging market", func() {
			highs := []float64{}
			lows := []float64{}
			closes := []float64{}

			for i := 0; i < 30; i++ {
				base := 105.0 + float64(i%2)*5 - 2.5
				highs = append(highs, base+3)
				lows = append(lows, base-3)
				closes = append(closes, base)
			}

			result, err := indicators.Stochastic(highs, lows, closes, 14, 3)

			Expect(err).NotTo(HaveOccurred())
			kValue, _ := result.K.Float64()
			Expect(kValue).To(BeNumerically("~", 50.0, 30.0))
		})
	})

	Describe("Standard parameters", func() {
		It("should work with standard 14,3 parameters", func() {
			highs := make([]float64, 50)
			lows := make([]float64, 50)
			closes := make([]float64, 50)

			for i := 0; i < 50; i++ {
				base := 100.0 + float64(i%10)
				highs[i] = base + 2
				lows[i] = base - 2
				closes[i] = base
			}

			result, err := indicators.Stochastic(highs, lows, closes, 14, 3)

			Expect(err).NotTo(HaveOccurred())
			kVal, _ := result.K.Float64()
			dVal, _ := result.D.Float64()
			Expect(kVal).To(BeNumerically(">=", 0.0))
			Expect(kVal).To(BeNumerically("<=", 100.0))
			Expect(dVal).To(BeNumerically(">=", 0.0))
			Expect(dVal).To(BeNumerically("<=", 100.0))
		})
	})
})
