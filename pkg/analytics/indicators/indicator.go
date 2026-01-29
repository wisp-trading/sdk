// Package indicators provides user-friendly methods for calculating technical indicators.
//
// This package wraps the low-level indicator calculations in pkg/analytics/indicators
// and handles all data fetching automatically. Users never need to manually extract
// klines or manage data - just call the indicator methods with an asset and period.
//
// Example usage:
//
//	btc := s.k.Asset("BTC")
//	rsi := s.k.Indicators.RSI(btc, 14)
//	sma := s.k.Indicators.SMA(btc, 50)
//	macd := s.k.Indicators.MACD(btc, 12, 26, 9)
//
// All indicator methods support optional configuration:
//
//	// Use specific exchange
//	rsi := s.k.Indicators.RSI(btc, 14, indicators.analytics.IndicatorOptions{
//	    Exchange: "binance",
//	})
//
//	// Use different timeframe
//	sma := s.k.Indicators.SMA(btc, 200, indicators.analytics.IndicatorOptions{
//	    Interval: "4h",
//	})
package indicators

import (
	"context"
	"fmt"
	"time"

	"github.com/wisp-trading/wisp/pkg/monitoring/profiling"
	"github.com/wisp-trading/wisp/pkg/types/connector"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
	"github.com/wisp-trading/wisp/pkg/types/wisp/analytics"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

const (
	// Fetch 2x period to ensure sufficient data for calculations
	dataMultiplier = 2
)

// IndicatorService provides user-friendly methods for technical indicators.
// All methods handle data fetching internally - users never manually extract klines.
type indicators struct {
	market analytics.Market
}

// NewIndicators creates a new IndicatorService
func NewIndicators(market analytics.Market) analytics.Indicators {
	return &indicators{
		market: market,
	}
}

// SMA calculates the Simple Moving Average for an asset.
//
// The Simple Moving Average is the average price over a specified number of periods.
// It's used to identify trends and smooth out price fluctuations.
//
// Parameters:
//   - asset: The asset to calculate SMA for (e.g., btc from s.k.Asset("BTC"))
//   - period: Number of periods to average (e.g., 20, 50, 200)
//   - opts: Optional exchange and interval configuration
//
// Returns the latest SMA value.
//
// Example:
//
//	btc := s.k.Asset("BTC")
//	sma50 := s.k.Indicators.SMA(btc, 50)  // 50-period SMA
//	sma200 := s.k.Indicators.SMA(btc, 200, analytics.analytics.IndicatorOptions{
//	    Interval: "4h",  // 4-hour timeframe
//	})
//
// Wisp automatically:
//   - Fetches price data from the exchange
//   - Calculates the moving average
//   - Returns the current value
func (s *indicators) SMA(ctx context.Context, asset portfolio.Asset, period int, opts ...analytics.IndicatorOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("SMA", time.Since(start))
		}
	}()

	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		return numerical.Zero(), err
	}

	return SMA(prices, period)
}

// EMA calculates the Exponential Moving Average for an asset.
//
// The Exponential Moving Average gives more weight to recent prices, making it more
// responsive to new information than SMA. Commonly used for trend identification
// and dynamic support/resistance levels.
//
// Parameters:
//   - asset: The asset to calculate EMA for
//   - period: Number of periods (e.g., 12, 20, 50, 200)
//   - opts: Optional exchange and interval configuration
//
// Returns the latest EMA value.
//
// Example:
//
//	eth := s.k.Asset("ETH")
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
func (s *indicators) EMA(ctx context.Context, asset portfolio.Asset, period int, opts ...analytics.IndicatorOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("EMA", time.Since(start))
		}
	}()

	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		return numerical.Zero(), err
	}

	return EMA(prices, period)
}

// RSI calculates the Relative Strength Index for an asset.
//
// RSI is a momentum oscillator that measures the speed and magnitude of price changes.
// Values range from 0 to 100:
//   - RSI > 70: Overbought (potential sell signal)
//   - RSI < 30: Oversold (potential buy signal)
//   - RSI = 50: Neutral
//
// Parameters:
//   - asset: The asset to calculate RSI for
//   - period: Number of periods (typically 14)
//   - opts: Optional exchange and interval configuration
//
// Returns the latest RSI value (0-100).
//
// Example:
//
//	btc := s.k.Asset("BTC")
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
func (s *indicators) RSI(ctx context.Context, asset portfolio.Asset, period int, opts ...analytics.IndicatorOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("RSI", time.Since(start))
		}
	}()

	prices, err := s.fetchClosePrices(asset, (period+1)*dataMultiplier, opts...)
	if err != nil {
		return numerical.Zero(), err
	}

	return RSI(prices, period)
}

// MACD calculates the Moving Average Convergence Divergence indicator.
//
// MACD is a trend-following momentum indicator that shows the relationship between
// two exponential moving averages. It consists of:
//   - MACD Line: Difference between fast and slow EMAs
//   - Signal Line: EMA of the MACD line
//   - Histogram: Difference between MACD and Signal lines
//
// Parameters:
//   - asset: The asset to calculate MACD for
//   - fastPeriod: Fast EMA period (typically 12)
//   - slowPeriod: Slow EMA period (typically 26)
//   - signalPeriod: Signal line period (typically 9)
//   - opts: Optional exchange and interval configuration
//
// Returns a MACDResult containing MACD, Signal, and Histogram values.
//
// Example:
//
//	btc := s.k.Asset("BTC")
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
func (s *indicators) MACD(ctx context.Context, asset portfolio.Asset, fastPeriod, slowPeriod, signalPeriod int, opts ...analytics.IndicatorOptions) (*analytics.MACDResult, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("MACD", time.Since(start))
		}
	}()

	requiredData := (slowPeriod + signalPeriod) * dataMultiplier
	prices, err := s.fetchClosePrices(asset, requiredData, opts...)
	if err != nil {
		return nil, err
	}

	macdResult, err := MACD(prices, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		return nil, err
	}

	return &macdResult, nil
}

// BollingerBands calculates Bollinger Bands for an asset.
//
// Bollinger Bands consist of three lines that envelope price action:
//   - Middle Band: Simple Moving Average
//   - Upper Band: Middle + (Standard Deviation × multiplier)
//   - Lower Band: Middle - (Standard Deviation × multiplier)
//
// They measure volatility and identify overbought/oversold conditions.
//
// Parameters:
//   - asset: The asset to calculate Bollinger Bands for
//   - period: Number of periods for SMA and std dev calculation (typically 20)
//   - stdDev: Standard deviation multiplier (typically 2.0)
//   - opts: Optional exchange and interval configuration
//
// Returns a BollingerBandsResult containing Upper, Middle, and Lower band values.
//
// Example:
//
//	btc := s.k.Asset("BTC")
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
func (s *indicators) BollingerBands(ctx context.Context, asset portfolio.Asset, period int, stdDev float64, opts ...analytics.IndicatorOptions) (*analytics.BollingerBandsResult, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("BollingerBands", time.Since(start))
		}
	}()

	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		return nil, err
	}

	bbResult, err := BollingerBands(prices, period, stdDev)
	if err != nil {
		return nil, err
	}

	return &bbResult, nil
}

// Stochastic calculates the Stochastic Oscillator for an asset.
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
//   - asset: The asset to calculate Stochastic for
//   - kPeriod: Lookback period for %K calculation (typically 14)
//   - dPeriod: Smoothing period for %D calculation (typically 3)
//   - opts: Optional exchange and interval configuration
//
// Returns a StochasticResult containing K and D values.
//
// Example:
//
//	eth := s.k.Asset("ETH")
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
func (s *indicators) Stochastic(ctx context.Context, asset portfolio.Asset, kPeriod, dPeriod int, opts ...analytics.IndicatorOptions) (*analytics.StochasticResult, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("Stochastic", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	requiredData := (kPeriod + dPeriod) * dataMultiplier
	klines := s.market.Klines(asset, exchange, interval, requiredData)
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	highs := make([]float64, len(klines))
	lows := make([]float64, len(klines))
	closes := make([]float64, len(klines))

	for i, kline := range klines {
		highs[i] = kline.High
		lows[i] = kline.Low
		closes[i] = kline.Close
	}

	stochResult, err := Stochastic(highs, lows, closes, kPeriod, dPeriod)
	if err != nil {
		return nil, err
	}

	return &stochResult, nil
}

// ATR calculates the Average True Range for an asset.
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
//   - asset: The asset to calculate ATR for
//   - period: Number of periods to average (typically 14)
//   - opts: Optional exchange and interval configuration
//
// Returns the latest ATR value in the asset's price units.
//
// Example:
//
//	btc := s.k.Asset("BTC")
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
func (s *indicators) ATR(ctx context.Context, asset portfolio.Asset, period int, opts ...analytics.IndicatorOptions) (numerical.Decimal, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("ATR", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return numerical.Zero(), fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	// Fetch klines (need high, low, close)
	requiredData := period * dataMultiplier
	klines := s.market.Klines(asset, exchange, interval, requiredData)
	if len(klines) == 0 {
		return numerical.Zero(), fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	// Extract high, low, close prices
	highs := make([]float64, len(klines))
	lows := make([]float64, len(klines))
	closes := make([]float64, len(klines))

	for i, kline := range klines {
		highs[i] = kline.High
		lows[i] = kline.Low
		closes[i] = kline.Close
	}

	return ATR(highs, lows, closes, period)
}

// fetchClosePrices is a helper that fetches klines and extracts close prices as float64
func (s *indicators) fetchClosePrices(asset portfolio.Asset, limit int, opts ...analytics.IndicatorOptions) ([]float64, error) {
	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.market.Klines(asset, exchange, interval, limit)
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	prices := make([]float64, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}

	return prices, nil
}

// getDefaultExchange returns the first available exchange for an asset
func (s *indicators) getDefaultExchange(asset portfolio.Asset) connector.ExchangeName {
	priceMap := s.market.Prices(context.Background(), asset)
	for exchange := range priceMap {
		return exchange
	}
	return ""
}

// parseOptions extracts options with defaults
func (s *indicators) parseOptions(opts ...analytics.IndicatorOptions) analytics.IndicatorOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.Interval == "" {
			options.Interval = analytics.DefaultInterval
		}
		return options
	}
	return analytics.IndicatorOptions{
		Interval: analytics.DefaultInterval,
	}
}
