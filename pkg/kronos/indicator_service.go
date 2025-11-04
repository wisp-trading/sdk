package kronos

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
	"github.com/shopspring/decimal"
)

const (
	// Default interval for kline fetching
	defaultInterval = "1h"
	// Fetch 2x period to ensure sufficient data for calculations
	dataMultiplier = 2
)

// IndicatorService provides user-friendly methods for technical indicators.
// All methods handle data fetching internally - users never manually extract klines.
type IndicatorService struct {
	store  store.Store
	logger logging.ApplicationLogger
}

// IndicatorOptions configures indicator calculations
type IndicatorOptions struct {
	Exchange connector.ExchangeName // Optional: defaults to first available exchange
	Interval string                 // Optional: defaults to 1h
}

// SMA calculates the Simple Moving Average for an asset.
// Returns the latest SMA value.
func (s *IndicatorService) SMA(asset portfolio.Asset, period int, opts ...IndicatorOptions) (decimal.Decimal, error) {
	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		s.logger.Warn("Failed to fetch prices for SMA", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	smaValues, err := indicators.SMA(prices, period)
	if err != nil {
		s.logger.Warn("Failed to calculate SMA", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	if len(smaValues) == 0 {
		return decimal.Zero, fmt.Errorf("no SMA values calculated")
	}

	// Return the latest value
	return smaValues[len(smaValues)-1], nil
}

// EMA calculates the Exponential Moving Average for an asset.
// Returns the latest EMA value.
func (s *IndicatorService) EMA(asset portfolio.Asset, period int, opts ...IndicatorOptions) (decimal.Decimal, error) {
	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		s.logger.Warn("Failed to fetch prices for EMA", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	emaValues, err := indicators.EMA(prices, period)
	if err != nil {
		s.logger.Warn("Failed to calculate EMA", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	if len(emaValues) == 0 {
		return decimal.Zero, fmt.Errorf("no EMA values calculated")
	}

	// Return the latest value
	return emaValues[len(emaValues)-1], nil
}

// RSI calculates the Relative Strength Index for an asset.
// Returns the latest RSI value (0-100).
func (s *IndicatorService) RSI(asset portfolio.Asset, period int, opts ...IndicatorOptions) (decimal.Decimal, error) {
	prices, err := s.fetchClosePrices(asset, (period+1)*dataMultiplier, opts...)
	if err != nil {
		s.logger.Warn("Failed to fetch prices for RSI", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	rsiValues, err := indicators.RSI(prices, period)
	if err != nil {
		s.logger.Warn("Failed to calculate RSI", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	if len(rsiValues) == 0 {
		return decimal.Zero, fmt.Errorf("no RSI values calculated")
	}

	// Return the latest value
	return rsiValues[len(rsiValues)-1], nil
}

// MACDResult holds the MACD indicator values
type MACDResult struct {
	MACD      decimal.Decimal
	Signal    decimal.Decimal
	Histogram decimal.Decimal
}

// MACD calculates the Moving Average Convergence Divergence for an asset.
// Returns the latest MACD, Signal, and Histogram values.
// Standard parameters: fastPeriod=12, slowPeriod=26, signalPeriod=9
func (s *IndicatorService) MACD(asset portfolio.Asset, fastPeriod, slowPeriod, signalPeriod int, opts ...IndicatorOptions) (*MACDResult, error) {
	// Need enough data for the slow period plus signal calculation
	requiredData := (slowPeriod + signalPeriod) * dataMultiplier
	prices, err := s.fetchClosePrices(asset, requiredData, opts...)
	if err != nil {
		s.logger.Warn("Failed to fetch prices for MACD", "asset", asset.Symbol(), "error", err)
		return nil, err
	}

	macdValues, err := indicators.MACD(prices, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		s.logger.Warn("Failed to calculate MACD", "asset", asset.Symbol(), "error", err)
		return nil, err
	}

	if len(macdValues) == 0 {
		return nil, fmt.Errorf("no MACD values calculated")
	}

	// Return the latest value
	latest := macdValues[len(macdValues)-1]
	return &MACDResult{
		MACD:      latest.MACD,
		Signal:    latest.Signal,
		Histogram: latest.Histogram,
	}, nil
}

// BollingerBandsResult holds the Bollinger Bands values
type BollingerBandsResult struct {
	Upper  decimal.Decimal
	Middle decimal.Decimal
	Lower  decimal.Decimal
}

// BollingerBands calculates Bollinger Bands for an asset.
// Returns the latest upper, middle, and lower band values.
// Standard parameters: period=20, stdDev=2.0
func (s *IndicatorService) BollingerBands(asset portfolio.Asset, period int, stdDev float64, opts ...IndicatorOptions) (*BollingerBandsResult, error) {
	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		s.logger.Warn("Failed to fetch prices for Bollinger Bands", "asset", asset.Symbol(), "period", period, "error", err)
		return nil, err
	}

	bbValues, err := indicators.BollingerBands(prices, period, stdDev)
	if err != nil {
		s.logger.Warn("Failed to calculate Bollinger Bands", "asset", asset.Symbol(), "period", period, "error", err)
		return nil, err
	}

	if len(bbValues) == 0 {
		return nil, fmt.Errorf("no Bollinger Bands values calculated")
	}

	// Return the latest value
	latest := bbValues[len(bbValues)-1]
	return &BollingerBandsResult{
		Upper:  latest.Upper,
		Middle: latest.Middle,
		Lower:  latest.Lower,
	}, nil
}

// StochasticResult holds the Stochastic Oscillator values
type StochasticResult struct {
	K decimal.Decimal // %K line (fast)
	D decimal.Decimal // %D line (slow)
}

// Stochastic calculates the Stochastic Oscillator for an asset.
// Returns the latest %K and %D values.
// Standard parameters: kPeriod=14, dPeriod=3
func (s *IndicatorService) Stochastic(asset portfolio.Asset, kPeriod, dPeriod int, opts ...IndicatorOptions) (*StochasticResult, error) {
	requiredData := (kPeriod + dPeriod) * dataMultiplier

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	// If no exchange specified, get first available
	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, requiredData)
	if len(klines) < kPeriod+dPeriod {
		return nil, fmt.Errorf("insufficient kline data for Stochastic: need %d, got %d", kPeriod+dPeriod, len(klines))
	}

	// Extract high, low, close prices
	highs := make([]decimal.Decimal, len(klines))
	lows := make([]decimal.Decimal, len(klines))
	closes := make([]decimal.Decimal, len(klines))

	for i, kline := range klines {
		highs[i] = kline.High
		lows[i] = kline.Low
		closes[i] = kline.Close
	}

	stochValues, err := indicators.Stochastic(highs, lows, closes, kPeriod, dPeriod)
	if err != nil {
		s.logger.Warn("Failed to calculate Stochastic", "asset", asset.Symbol(), "error", err)
		return nil, err
	}

	if len(stochValues) == 0 {
		return nil, fmt.Errorf("no Stochastic values calculated")
	}

	// Return the latest value
	latest := stochValues[len(stochValues)-1]
	return &StochasticResult{
		K: latest.K,
		D: latest.D,
	}, nil
}

// ATR calculates the Average True Range for an asset.
// Returns the latest ATR value.
// Standard parameter: period=14
func (s *IndicatorService) ATR(asset portfolio.Asset, period int, opts ...IndicatorOptions) (decimal.Decimal, error) {
	requiredData := (period + 1) * dataMultiplier

	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	// If no exchange specified, get first available
	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return decimal.Zero, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, requiredData)
	if len(klines) < period+1 {
		return decimal.Zero, fmt.Errorf("insufficient kline data for ATR: need %d, got %d", period+1, len(klines))
	}

	// Extract high, low, close prices
	highs := make([]decimal.Decimal, len(klines))
	lows := make([]decimal.Decimal, len(klines))
	closes := make([]decimal.Decimal, len(klines))

	for i, kline := range klines {
		highs[i] = kline.High
		lows[i] = kline.Low
		closes[i] = kline.Close
	}

	atrValues, err := indicators.ATR(highs, lows, closes, period)
	if err != nil {
		s.logger.Warn("Failed to calculate ATR", "asset", asset.Symbol(), "period", period, "error", err)
		return decimal.Zero, err
	}

	if len(atrValues) == 0 {
		return decimal.Zero, fmt.Errorf("no ATR values calculated")
	}

	// Return the latest value
	return atrValues[len(atrValues)-1], nil
}

// fetchClosePrices is a helper that fetches klines and extracts close prices
func (s *IndicatorService) fetchClosePrices(asset portfolio.Asset, limit int, opts ...IndicatorOptions) ([]decimal.Decimal, error) {
	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	// If no exchange specified, get first available
	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	klines := s.store.GetKlines(asset, exchange, interval, limit)
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
	}

	// Extract close prices
	prices := make([]decimal.Decimal, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}

	return prices, nil
}

// getDefaultExchange returns the first available exchange for an asset
func (s *IndicatorService) getDefaultExchange(asset portfolio.Asset) connector.ExchangeName {
	// Try to get price data to find available exchanges
	priceMap := s.store.GetAssetPrices(asset)
	for exchange := range priceMap {
		return exchange
	}
	return ""
}

// parseOptions extracts options with defaults
func (s *IndicatorService) parseOptions(opts ...IndicatorOptions) IndicatorOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.Interval == "" {
			options.Interval = defaultInterval
		}
		return options
	}
	return IndicatorOptions{
		Interval: defaultInterval,
	}
}
