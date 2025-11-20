package analytics

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/shopspring/decimal"
)

type Analytics interface {
	Volatility(asset portfolio.Asset, period int, opts ...AnalyticsOptions) (decimal.Decimal, error)
	Trend(asset portfolio.Asset, period int, opts ...AnalyticsOptions) (*TrendResult, error)
	VolumeAnalysis(asset portfolio.Asset, period int, opts ...AnalyticsOptions) (*VolumeAnalysis, error)
	GetPriceChange(asset portfolio.Asset, period int, opts ...AnalyticsOptions) (*PriceChange, error)
}

// AnalyticsOptions configures analytics calculations
type AnalyticsOptions struct {
	Exchange connector.ExchangeName
	Interval string
}

// TrendDirection represents the trend direction
type TrendDirection string

const (
	TrendBullish TrendDirection = "bullish"
	TrendBearish TrendDirection = "bearish"
	TrendNeutral TrendDirection = "neutral"
)

// TrendResult holds trend analysis results
type TrendResult struct {
	Direction TrendDirection
	Strength  decimal.Decimal // 0-100, higher means stronger trend
	Slope     decimal.Decimal // Linear regression slope
}

// PriceChange calculates the price change over a period.
type PriceChange struct {
	StartPrice        decimal.Decimal
	EndPrice          decimal.Decimal
	Change            decimal.Decimal // Absolute change
	ChangePercent     decimal.Decimal // Percentage change
	HighPrice         decimal.Decimal // Highest price in period
	LowPrice          decimal.Decimal // Lowest price in period
	PriceRange        decimal.Decimal // High - Low
	PriceRangePercent decimal.Decimal // Range as % of start price
}

// VolumeAnalysis holds volume analysis results
type VolumeAnalysis struct {
	CurrentVolume decimal.Decimal
	AverageVolume decimal.Decimal
	VolumeRatio   decimal.Decimal // Current / Average
	IsVolumeSpike bool            // True if current volume > 2x average
	VolumeTrend   TrendDirection  // Increasing, decreasing, or neutral
}
