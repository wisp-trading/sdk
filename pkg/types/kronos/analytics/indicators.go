package analytics

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

type Indicators interface {
	ATR(ctx context.Context, asset portfolio.Asset, period int, opts ...IndicatorOptions) (numerical.Decimal, error)
	SMA(ctx context.Context, asset portfolio.Asset, period int, opts ...IndicatorOptions) (numerical.Decimal, error)
	EMA(ctx context.Context, asset portfolio.Asset, period int, opts ...IndicatorOptions) (numerical.Decimal, error)
	RSI(ctx context.Context, asset portfolio.Asset, period int, opts ...IndicatorOptions) (numerical.Decimal, error)
	MACD(ctx context.Context, asset portfolio.Asset, fastPeriod, slowPeriod, signalPeriod int, opts ...IndicatorOptions) (*MACDResult, error)
	BollingerBands(ctx context.Context, asset portfolio.Asset, period int, stdDev float64, opts ...IndicatorOptions) (*BollingerBandsResult, error)
	Stochastic(ctx context.Context, asset portfolio.Asset, kPeriod, dPeriod int, opts ...IndicatorOptions) (*StochasticResult, error)
}

// IndicatorOptions configures indicator calculations.
// All fields are optional. If not specified, Kronos uses sensible defaults:
// - Exchange: First available exchange with data for the asset
// - Interval: 1h (hourly candles)
type IndicatorOptions struct {
	// Exchange specifies which exchange to fetch data from.
	// If empty, Kronos automatically selects the first available exchange.
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
