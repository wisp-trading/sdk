package analytics

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Analytics provides higher-level market analysis built on top of kline data.
//
// All methods are pure functions — callers are responsible for fetching klines
// (e.g. via wisp.Spot().Klines(...) or wisp.Perp().Klines(...)) and passing them in.
// This keeps data fetching explicit and avoids hidden store access inside analytics calls.
//
// Example:
//
//	klines := wisp.Spot().Klines(exchange, btc, "1h", 60)
//	trend, _    := wisp.Analytics().Trend(klines)
//	vol, _      := wisp.Analytics().Volatility(klines, "1h")
//	change, _   := wisp.Analytics().GetPriceChange(klines)
//
// TODO: Implement equivalent analytics directly on wisp.Spot() and wisp.Perp() domain objects
// so strategies can call e.g. wisp.Spot().Trend(exchange, btc, "1h", 60) without managing
// kline fetching themselves. Shared computation logic should live under markets/base.
type Analytics interface {
	// Volatility calculates annualised volatility (std dev of returns) from the provided klines.
	// interval is used to derive the annualisation factor (e.g. "1h", "4h", "1d").
	Volatility(klines []connector.Kline, interval string) (numerical.Decimal, error)

	// Trend analyses price trend via linear regression over the provided klines.
	Trend(klines []connector.Kline) (*TrendResult, error)

	// VolumeAnalysis detects volume patterns and spikes from the provided klines.
	VolumeAnalysis(klines []connector.Kline) (*VolumeAnalysis, error)

	// GetPriceChange calculates price statistics over the provided klines.
	GetPriceChange(klines []connector.Kline) (*PriceChange, error)
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
