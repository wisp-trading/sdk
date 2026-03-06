// Package indicators provides technical indicator calculations.
// All methods are pure functions — callers are responsible for fetching klines
// (e.g. via wisp.Spot().Klines(...) or wisp.Perp().Klines(...)) and passing them in.
package indicators

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// indicators is the concrete implementation of analytics.Indicators.
type indicators struct{}

// NewIndicators creates a new Indicators service.
func NewIndicators() analytics.Indicators {
	return &indicators{}
}

// SMA calculates the Simple Moving Average from the provided klines.
//
// The Simple Moving Average is the average price over a specified number of periods.
// It's used to identify trends and smooth out price fluctuations.
//
// Parameters:
//   - klines: The kline data to calculate SMA from
//   - period: Number of periods to average (e.g., 20, 50, 200)
//
// Returns the latest SMA value.
//
// Example:
//
//	btc := s.k.Pair("BTC")
//	sma50 := s.k.Indicators.SMA(btc, 50)  // 50-period SMA
//	sma200 := s.k.Indicators.SMA(btc, 200, analytics.analytics.IndicatorOptions{
//	    Interval: "4h",  // 4-hour timeframe
//	})
//
// Wisp automatically:
//   - Fetches price data from the exchange
//   - Calculates the moving average
//   - Returns the current value
func (s *indicators) SMA(klines []connector.Kline, period int) (numerical.Decimal, error) {
	return SMA(extractClose(klines), period)
}

// EMA calculates the Exponential Moving Average from the provided klines.
//
// The Exponential Moving Average gives more weight to recent prices, making it more
// responsive to new information than SMA. Commonly used for trend identification
// and dynamic support/resistance levels.
//
// Parameters:
//   - klines: The kline data to calculate EMA from
//   - period: Number of periods (e.g., 12, 20, 50, 200)
//
// Returns the latest EMA value.
//
// Example:
//
//	eth := s.k.Pair("ETH")
//	ema20 := s.k.Indicators.EMA(eth, 20)   // 20-period EMA
//	ema50 := s.k.Indicators.EMA(eth, 50)   // 50-period EMA
//
//	// Check if price is above EMA (uptrend)
//	price := s.k.Market.Price(eth)
//	if price.GreaterThan(ema50) {
//	    // Uptrend detected
//	}
//
// Wisp automatically:
//   - Fetches historical price data
//   - Applies exponential weighting
//   - Returns the current EMA value
func (s *indicators) EMA(klines []connector.Kline, period int) (numerical.Decimal, error) {
	return EMA(extractClose(klines), period)
}

// RSI calculates the Relative Strength Index from the provided klines.
//
// RSI is a momentum oscillator that measures the speed and magnitude of price changes.
// Values range from 0 to 100:
//   - RSI > 70: Overbought (potential sell signal)
//   - RSI < 30: Oversold (potential buy signal)
//   - RSI = 50: Neutral
//
// Parameters:
//   - klines: The kline data to calculate RSI from
//   - period: Number of periods (typically 14)
//
// Returns the latest RSI value (0-100).
//
// Example:
//
//	btc := s.k.Pair("BTC")
//	rsi := s.k.Indicators.RSI(btc, 14)  // Standard 14-period RSI
//
//	if rsi.LessThan(numerical.NewFromInt(30)) {
//	    // Oversold - potential buy signal
//	    return s.Signal().Buy(btc).Build()
//	}
//
//	if rsi.GreaterThan(numerical.NewFromInt(70)) {
//	    // Overbought - potential sell signal
//	    return s.Signal().Sell(btc).Build()
//	}
//
// Wisp automatically:
//   - Fetches price history
//   - Calculates gains and losses
//   - Computes the RSI value
func (s *indicators) RSI(klines []connector.Kline, period int) (numerical.Decimal, error) {
	return RSI(extractClose(klines), period)
}

// MACD calculates the Moving Average Convergence Divergence from the provided klines.
//
// MACD is a trend-following momentum indicator that shows the relationship between
// two exponential moving averages. It consists of:
//   - MACD Line: Difference between fast and slow EMAs
//   - Signal Line: EMA of the MACD line
//   - Histogram: Difference between MACD and Signal lines
//
// Parameters:
//   - klines: The kline data to calculate MACD from
//   - fastPeriod: Fast EMA period (typically 12)
//   - slowPeriod: Slow EMA period (typically 26)
//   - signalPeriod: Signal line period (typically 9)
//
// Returns a MACDResult containing MACD, Signal, and Histogram values.
//
// Example:
//
//	btc := s.k.Pair("BTC")
//	macd := s.k.Indicators.MACD(btc, 12, 26, 9)  // Standard settings
//
//	// Bullish crossover: MACD crosses above signal
//	if macd.MACD.GreaterThan(macd.Signal) {
//	    return s.Signal().Buy(btc).Reason("MACD bullish crossover").Build()
//	}
//
//	// Bearish crossover: MACD crosses below signal
//	if macd.MACD.LessThan(macd.Signal) {
//	    return s.Signal().Sell(btc).Reason("MACD bearish crossover").Build()
//	}
//
// Interpretation:
//   - MACD > Signal: Bullish momentum
//   - MACD < Signal: Bearish momentum
//   - Histogram growing: Momentum strengthening
//   - Histogram shrinking: Momentum weakening
func (s *indicators) MACD(klines []connector.Kline, fastPeriod, slowPeriod, signalPeriod int) (*analytics.MACDResult, error) {
	result, err := MACD(extractClose(klines), fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BollingerBands calculates Bollinger Bands from the provided klines.
//
// Bollinger Bands consist of three lines that envelope price action:
//   - Middle Band: Simple Moving Average
//   - Upper Band: Middle + (Standard Deviation × multiplier)
//   - Lower Band: Middle - (Standard Deviation × multiplier)
//
// They measure volatility and identify overbought/oversold conditions.
//
// Parameters:
//   - klines: The kline data to calculate Bollinger Bands from
//   - period: Number of periods for SMA and std dev calculation (typically 20)
//   - stdDev: Standard deviation multiplier (typically 2.0)
//
// Returns a BollingerBandsResult containing Upper, Middle, and Lower band values.
//
// Example:
//
//	btc := s.k.Pair("BTC")
//	bb := s.k.Indicators.BollingerBands(btc, 20, 2.0)  // 20-period, 2 std dev
//	price := s.k.Market.Price(btc)
//
//	// Price touching lower band - potential buy
//	if price.LessThan(bb.Lower) {
//	    return s.Signal().Buy(btc).Reason("Price below lower Bollinger Band").Build()
//	}
//
//	// Price touching upper band - potential sell
//	if price.GreaterThan(bb.Upper) {
//	    return s.Signal().Sell(btc).Reason("Price above upper Bollinger Band").Build()
//	}
//
// Interpretation:
//   - Price near upper band: Overbought
//   - Price near lower band: Oversold
//   - Bands narrowing: Low volatility (potential breakout)
//   - Bands widening: High volatility
func (s *indicators) BollingerBands(klines []connector.Kline, period int, stdDev float64) (*analytics.BollingerBandsResult, error) {
	result, err := BollingerBands(extractClose(klines), period, stdDev)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Stochastic calculates the Stochastic Oscillator from the provided klines.
//
// The Stochastic Oscillator compares a closing price to its price range over a period.
// It consists of two lines:
//   - %K (Fast): Current position within the price range (0-100)
//   - %D (Slow): Moving average of %K, providing a smoother signal
//
// Values range from 0 to 100:
//   - > 80: Overbought
//   - < 20: Oversold
//
// Parameters:
//   - klines: The kline data to calculate Stochastic from
//   - kPeriod: Lookback period for %K calculation (typically 14)
//   - dPeriod: Smoothing period for %D calculation (typically 3)
//
// Returns a StochasticResult containing K and D values.
//
// Example:
//
//	eth := s.k.Pair("ETH")
//	stoch := s.k.Indicators.Stochastic(eth, 14, 3)  // Standard 14,3 settings
//
//	// Both lines oversold - strong buy signal
//	if stoch.K.LessThan(numerical.NewFromInt(20)) &&
//	   stoch.D.LessThan(numerical.NewFromInt(20)) {
//	    return s.Signal().Buy(eth).Reason("Stochastic oversold").Build()
//	}
//
//	// Both lines overbought - strong sell signal
//	if stoch.K.GreaterThan(numerical.NewFromInt(80)) &&
//	   stoch.D.GreaterThan(numerical.NewFromInt(80)) {
//	    return s.Signal().Sell(eth).Reason("Stochastic overbought").Build()
//	}
//
// Interpretation:
//   - %K crosses above %D: Bullish signal
//   - %K crosses below %D: Bearish signal
//   - Both in oversold zone: Potential reversal up
//   - Both in overbought zone: Potential reversal down
func (s *indicators) Stochastic(klines []connector.Kline, kPeriod, dPeriod int) (*analytics.StochasticResult, error) {
	highs, lows, closes := extractHLC(klines)
	result, err := Stochastic(highs, lows, closes, kPeriod, dPeriod)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ATR calculates the Average True Range from the provided klines.
//
// ATR measures market volatility by calculating the average range between
// high and low prices over a period. Higher ATR indicates higher volatility.
//
// ATR is commonly used for:
//   - Setting stop losses (e.g., stop at 2× ATR below entry)
//   - Position sizing (reduce size in high volatility)
//   - Identifying breakouts (ATR expanding)
//
// Parameters:
//   - klines: The kline data to calculate ATR from
//   - period: Number of periods to average (typically 14)
//
// Returns the latest ATR value in the asset's price units.
//
// Example:
//
//	btc := s.k.Pair("BTC")
//	atr := s.k.Indicators.ATR(btc, 14)  // 14-period ATR
//	price := s.k.Market.Price(btc)
//
//	// Set stop loss at 2× ATR below entry
//	stopLoss := price.Sub(atr.Mul(numerical.NewFromInt(2)))
//
//	// Check if volatility is high
//	avgPrice := s.k.Indicators.SMA(btc, 20)
//	atrPercent := atr.Div(avgPrice).Mul(numerical.NewFromInt(100))
//	if atrPercent.GreaterThan(numerical.NewFromInt(5)) {
//	    // High volatility - reduce position size
//	}
//
// Interpretation:
//   - High ATR: High volatility, larger price swings
//   - Low ATR: Low volatility, price consolidation
//   - Rising ATR: Volatility increasing
//   - Falling ATR: Volatility decreasing
func (s *indicators) ATR(klines []connector.Kline, period int) (numerical.Decimal, error) {
	highs, lows, closes := extractHLC(klines)
	return ATR(highs, lows, closes, period)
}

// ============================================================
// Internal helpers
// ============================================================

func extractClose(klines []connector.Kline) []float64 {
	prices := make([]float64, len(klines))
	for i, k := range klines {
		prices[i] = k.Close
	}
	return prices
}

func extractHLC(klines []connector.Kline) (highs, lows, closes []float64) {
	highs = make([]float64, len(klines))
	lows = make([]float64, len(klines))
	closes = make([]float64, len(klines))
	for i, k := range klines {
		highs[i] = k.High
		lows[i] = k.Low
		closes[i] = k.Close
	}
	return
}
