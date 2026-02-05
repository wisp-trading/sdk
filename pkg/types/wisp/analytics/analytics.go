package analytics

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type Analytics interface {
	Volatility(ctx context.Context, asset portfolio.Pair, period int, opts ...AnalyticsOptions) (numerical.Decimal, error)
	Trend(ctx context.Context, asset portfolio.Pair, period int, opts ...AnalyticsOptions) (*TrendResult, error)
	VolumeAnalysis(ctx context.Context, asset portfolio.Pair, period int, opts ...AnalyticsOptions) (*VolumeAnalysis, error)
	GetPriceChange(ctx context.Context, asset portfolio.Pair, period int, opts ...AnalyticsOptions) (*PriceChange, error)
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
	Strength  numerical.Decimal // 0-100, higher means stronger trend
	Slope     numerical.Decimal // Linear regression slope
}

// PriceChange calculates the price change over a period.
type PriceChange struct {
	StartPrice        numerical.Decimal
	EndPrice          numerical.Decimal
	Change            numerical.Decimal // Absolute change
	ChangePercent     numerical.Decimal // Percentage change
	HighPrice         numerical.Decimal // Highest price in period
	LowPrice          numerical.Decimal // Lowest price in period
	PriceRange        numerical.Decimal // High - Low
	PriceRangePercent numerical.Decimal // Range as % of start price
}

// VolumeAnalysis holds volume analysis results
type VolumeAnalysis struct {
	CurrentVolume numerical.Decimal
	AverageVolume numerical.Decimal
	VolumeRatio   numerical.Decimal // Current / Average
	IsVolumeSpike bool              // True if current volume > 2x average
	VolumeTrend   TrendDirection    // Increasing, decreasing, or neutral
}
