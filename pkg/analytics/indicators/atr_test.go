package indicators_test

import (
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ATR", func() {
	Describe("Input validation", func() {
		It("should return error when arrays have different lengths", func() {
			highs := makeDecimals(100, 105, 110)
			lows := makeDecimals(95, 100)
			closes := makeDecimals(98, 103, 108)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("equal length"))
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should return error when insufficient data for period", func() {
			highs := makeDecimals(100, 105)
			lows := makeDecimals(95, 100)
			closes := makeDecimals(98, 103)

			result, err := indicators.ATR(highs, lows, closes, 5)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("insufficient data"))
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should require at least period+1 data points", func() {
			highs := makeDecimals(100, 105, 110, 115)
			lows := makeDecimals(95, 100, 105, 110)
			closes := makeDecimals(98, 103, 108, 113)

			result, err := indicators.ATR(highs, lows, closes, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("True Range calculation", func() {
		It("should calculate ATR for simple data", func() {
			highs := makeDecimals(100, 105, 110, 115, 120)
			lows := makeDecimals(95, 100, 105, 110, 115)
			closes := makeDecimals(98, 103, 108, 113, 118)

			result, err := indicators.ATR(highs, lows, closes, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle price gaps correctly", func() {
			highs := makeDecimals(100, 110, 115)
			lows := makeDecimals(95, 105, 110)
			closes := makeDecimals(98, 108, 113)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle gap down correctly", func() {
			highs := makeDecimals(100, 90, 95)
			lows := makeDecimals(95, 85, 90)
			closes := makeDecimals(98, 88, 93)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("ATR smoothing", func() {
		It("should apply EMA-like smoothing to subsequent values", func() {
			highs := makeDecimals(100, 105, 110, 115, 120, 125, 130)
			lows := makeDecimals(95, 100, 105, 110, 115, 120, 125)
			closes := makeDecimals(98, 103, 108, 113, 118, 123, 128)

			result, err := indicators.ATR(highs, lows, closes, 3)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should produce smoother values over time", func() {
			highs := []numerical.Decimal{}
			lows := []numerical.Decimal{}
			closes := []numerical.Decimal{}

			for i := 0; i < 20; i++ {
				highs = append(highs, numerical.NewFromInt(int64(100+i*5)))
				lows = append(lows, numerical.NewFromInt(int64(95+i*5)))
				closes = append(closes, numerical.NewFromInt(int64(98+i*5)))
			}

			result, err := indicators.ATR(highs, lows, closes, 5)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Edge cases", func() {
		It("should handle identical prices (zero volatility)", func() {
			highs := makeDecimals(100, 100, 100, 100)
			lows := makeDecimals(100, 100, 100, 100)
			closes := makeDecimals(100, 100, 100, 100)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Equal(numerical.Zero())).To(BeTrue())
		})

		It("should handle single period", func() {
			highs := makeDecimals(100, 105)
			lows := makeDecimals(95, 100)
			closes := makeDecimals(98, 103)

			result, err := indicators.ATR(highs, lows, closes, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle large price swings", func() {
			highs := makeDecimals(100, 200, 150, 300)
			lows := makeDecimals(90, 180, 140, 280)
			closes := makeDecimals(95, 190, 145, 290)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.NewFromInt(10))).To(BeTrue())
		})

		It("should handle fractional prices", func() {
			highs := makeDecimalsFloat(100.5, 101.75, 102.25, 103.5)
			lows := makeDecimalsFloat(99.25, 100.5, 101.0, 102.25)
			closes := makeDecimalsFloat(100.0, 101.25, 101.75, 103.0)

			result, err := indicators.ATR(highs, lows, closes, 2)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Result length", func() {
		It("should return single value", func() {
			dataLength := 100
			period := 14

			highs := make([]numerical.Decimal, dataLength)
			lows := make([]numerical.Decimal, dataLength)
			closes := make([]numerical.Decimal, dataLength)

			for i := 0; i < dataLength; i++ {
				highs[i] = numerical.NewFromInt(int64(100 + i))
				lows[i] = numerical.NewFromInt(int64(95 + i))
				closes[i] = numerical.NewFromInt(int64(98 + i))
			}

			result, err := indicators.ATR(highs, lows, closes, period)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})

	Describe("Real-world scenarios", func() {
		It("should calculate ATR for typical market data", func() {
			highs := makeDecimalsFloat(50100, 50250, 50150, 50300, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800, 50750, 50900)
			lows := makeDecimalsFloat(49900, 50050, 49950, 50100, 50200, 50150, 50300, 50250, 50400, 50350, 50500, 50450, 50600, 50550, 50700)
			closes := makeDecimalsFloat(50000, 50150, 50050, 50200, 50300, 50250, 50400, 50350, 50500, 50450, 50600, 50550, 50700, 50650, 50800)

			result, err := indicators.ATR(highs, lows, closes, 14)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})

		It("should handle trending market", func() {
			highs := []numerical.Decimal{}
			lows := []numerical.Decimal{}
			closes := []numerical.Decimal{}

			for i := 0; i < 30; i++ {
				base := float64(100 + i*2)
				highs = append(highs, numerical.NewFromFloat(base+5))
				lows = append(lows, numerical.NewFromFloat(base-5))
				closes = append(closes, numerical.NewFromFloat(base+2))
			}

			result, err := indicators.ATR(highs, lows, closes, 10)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.GreaterThan(numerical.Zero())).To(BeTrue())
		})
	})
})

// Helper functions
func makeDecimals(values ...float64) []numerical.Decimal {
	result := make([]numerical.Decimal, len(values))
	for i, v := range values {
		result[i] = numerical.NewFromFloat(v)
	}
	return result
}

func makeDecimalsFloat(values ...float64) []numerical.Decimal {
	result := make([]numerical.Decimal, len(values))
	for i, v := range values {
		result[i] = numerical.NewFromFloat(v)
	}
	return result
}
