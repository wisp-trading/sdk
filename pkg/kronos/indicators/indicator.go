package indicators

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/types/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/stores/market"
	"github.com/shopspring/decimal"
)

const (
	// Fetch 2x period to ensure sufficient data for calculations
	dataMultiplier = 2
)

// IndicatorService provides user-friendly methods for technical indicators.
// All methods handle data fetching internally - users never manually extract klines.
type IndicatorService struct {
	store market.MarketData
}

// NewIndicatorService creates a new IndicatorService
func NewIndicatorService(store market.MarketData) *IndicatorService {
	return &IndicatorService{
		store: store,
	}
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
		return decimal.Zero, err
	}

	smaValues, err := indicators.SMA(prices, period)
	if err != nil {
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
		return decimal.Zero, err
	}

	emaValues, err := indicators.EMA(prices, period)
	if err != nil {
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
		return decimal.Zero, err
	}

	rsiValues, err := indicators.RSI(prices, period)
	if err != nil {
		return decimal.Zero, err
	}

	if len(rsiValues) == 0 {
		return decimal.Zero, fmt.Errorf("no RSI values calculated")
	}

	// Return the latest value
	return rsiValues[len(rsiValues)-1], nil
}

// MACD calculates the Moving Average Convergence Divergence indicator.
// Returns the latest MACD, signal, and histogram values.
func (s *IndicatorService) MACD(asset portfolio.Asset, fastPeriod, slowPeriod, signalPeriod int, opts ...IndicatorOptions) (*indicators.MACDResult, error) {
	// Need enough data for slow period + signal period
	requiredData := (slowPeriod + signalPeriod) * dataMultiplier
	prices, err := s.fetchClosePrices(asset, requiredData, opts...)
	if err != nil {
		return nil, err
	}

	macdResults, err := indicators.MACD(prices, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		return nil, err
	}

	if len(macdResults) == 0 {
		return nil, fmt.Errorf("no MACD values calculated")
	}

	// Return the latest values
	return &macdResults[len(macdResults)-1], nil
}

// BollingerBands calculates Bollinger Bands for an asset.
// Returns the latest upper, middle, and lower band values.
func (s *IndicatorService) BollingerBands(asset portfolio.Asset, period int, stdDev float64, opts ...IndicatorOptions) (*indicators.BollingerBandsResult, error) {
	prices, err := s.fetchClosePrices(asset, period*dataMultiplier, opts...)
	if err != nil {
		return nil, err
	}

	bbResults, err := indicators.BollingerBands(prices, period, stdDev)
	if err != nil {
		return nil, err
	}

	if len(bbResults) == 0 {
		return nil, fmt.Errorf("no Bollinger Bands values calculated")
	}

	// Return the latest values
	return &bbResults[len(bbResults)-1], nil
}

// Stochastic calculates the Stochastic oscillator for an asset.
// Returns the latest %K and %D values.
func (s *IndicatorService) Stochastic(asset portfolio.Asset, kPeriod, dPeriod int, opts ...IndicatorOptions) (*indicators.StochasticResult, error) {
	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return nil, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	// Fetch klines (need high, low, close)
	requiredData := (kPeriod + dPeriod) * dataMultiplier
	klines := s.store.GetKlines(asset, exchange, interval, requiredData)
	if len(klines) == 0 {
		return nil, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
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

	stochResults, err := indicators.Stochastic(highs, lows, closes, kPeriod, dPeriod)
	if err != nil {
		return nil, err
	}

	if len(stochResults) == 0 {
		return nil, fmt.Errorf("no Stochastic values calculated")
	}

	// Return the latest values
	return &stochResults[len(stochResults)-1], nil
}

// ATR calculates the Average True Range for an asset.
// Returns the latest ATR value.
func (s *IndicatorService) ATR(asset portfolio.Asset, period int, opts ...IndicatorOptions) (decimal.Decimal, error) {
	options := s.parseOptions(opts...)
	exchange := options.Exchange
	interval := options.Interval

	if exchange == "" {
		exchange = s.getDefaultExchange(asset)
		if exchange == "" {
			return decimal.Zero, fmt.Errorf("no exchange data available for asset %s", asset.Symbol())
		}
	}

	// Fetch klines (need high, low, close)
	requiredData := period * dataMultiplier
	klines := s.store.GetKlines(asset, exchange, interval, requiredData)
	if len(klines) == 0 {
		return decimal.Zero, fmt.Errorf("no kline data available for asset %s on exchange %s", asset.Symbol(), exchange)
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

	prices := make([]decimal.Decimal, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}

	return prices, nil
}

// getDefaultExchange returns the first available exchange for an asset
func (s *IndicatorService) getDefaultExchange(asset portfolio.Asset) connector.ExchangeName {
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
			options.Interval = analytics.DefaultInterval
		}
		return options
	}
	return IndicatorOptions{
		Interval: analytics.DefaultInterval,
	}
}
