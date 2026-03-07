package analytics

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Indicators provides technical indicator calculations.
//
// Two usage patterns are supported:
//
// 1. Fetch-and-compute — supply an asset, exchange, and period. The service fetches
// klines from the appropriate store and runs the calculation. Use this for simple
// single-indicator strategies.
//
// 2. Compute-only () — supply pre-fetched []connector.Kline directly.
// The service runs the calculation as a pure function with no store access. Use
// this when multiple indicators share the same kline window in a single tick,
// avoiding redundant store reads.
//
// Example — compute-only pattern:
//
//	klines := wisp.Market().Spot().GetKlines(btc, "binance", analytics.Interval1Hour, 60)
//	rsi, _  := wisp.Indicators().RSI(klines, 14)
//	macd, _ := wisp.Indicators().MACD(klines, 12, 26, 9)
//	bb, _   := wisp.Indicators().BollingerBands(klines, 20, 2.0)
type Indicators interface {
	ATR(klines []connector.Kline, period int) (numerical.Decimal, error)
	SMA(klines []connector.Kline, period int) (numerical.Decimal, error)
	EMA(klines []connector.Kline, period int) (numerical.Decimal, error)
	RSI(klines []connector.Kline, period int) (numerical.Decimal, error)
	MACD(klines []connector.Kline, fastPeriod, slowPeriod, signalPeriod int) (*MACDResult, error)
	BollingerBands(klines []connector.Kline, period int, stdDev float64) (*BollingerBandsResult, error)
	Stochastic(klines []connector.Kline, kPeriod, dPeriod int) (*StochasticResult, error)
}

// IndicatorOptions configures fetch-and-compute indicator calls.
// Exchange is required — the caller must know which market and exchange the data
// should come from. Use wisp.Market().Spot() or wisp.Market().Perp() to discover
// available exchanges for an asset.
type IndicatorOptions struct {
	// Exchange specifies which exchange to fetch kline data from. Required.
	Exchange connector.ExchangeName

	// Interval specifies the timeframe for calculations (e.g., "1h", "4h", "1d").
	// If empty, defaults to "1h".
	Interval string
}

type BollingerBandsResult struct {
	Upper  numerical.Decimal
	Middle numerical.Decimal
	Lower  numerical.Decimal
}

type MACDResult struct {
	MACD      numerical.Decimal
	Signal    numerical.Decimal
	Histogram numerical.Decimal
}

// StochasticResult represents the output of the Stochastic Oscillator calculation.
// The Stochastic Oscillator is a momentum indicator that compares a security's closing price
// to its price range over a given time period.
type StochasticResult struct {
	// K represents the %K line (fast stochastic), which measures the current close
	// relative to the high-low range over the specified period. Values range from 0 to 100.
	K numerical.Decimal

	// D represents the %D line (slow stochastic), which is the simple moving average
	// of the %K line over the specified period. Values range from 0 to 100.
	// The %D line provides a smoothed signal line.
	D numerical.Decimal
}
