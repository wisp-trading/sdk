package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/wisp-trading/wisp/pkg/monitoring/profiling"
	"github.com/wisp-trading/wisp/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/wisp/pkg/types/wisp/analytics"
	"github.com/wisp-trading/wisp/pkg/types/wisp/numerical"
)

// VolumeAnalysis detects volume patterns and spikes.
func (s *analytics) VolumeAnalysis(ctx context.Context, asset portfolio.Asset, period int, opts ...analyticsTypes.AnalyticsOptions) (*analyticsTypes.VolumeAnalysis, error) {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("VolumeAnalysis", time.Since(start))
		}
	}()

	options := s.parseOptions(opts...)

	klines, err := s.fetchKlines(ctx, asset, period+1, options)
	if err != nil {
		return nil, err
	}

	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for volume analysis")
	}

	// Get current volume (latest kline)
	currentVolume := klines[len(klines)-1].Volume

	// Calculate average volume over period
	totalVolume := 0.0
	for i := 0; i < len(klines)-1; i++ {
		totalVolume += klines[i].Volume
	}
	avgVolume := totalVolume / float64(len(klines)-1)

	// Calculate volume ratio
	var volumeRatio float64
	if avgVolume == 0 {
		volumeRatio = 0
	} else {
		volumeRatio = currentVolume / avgVolume
	}

	// Detect spike (current > 2x average)
	isSpike := volumeRatio > 2.0

	// Determine volume trend
	// Compare first half average to second half average
	midPoint := len(klines) / 2
	firstHalfVolume := 0.0
	secondHalfVolume := 0.0

	for i := 0; i < midPoint; i++ {
		firstHalfVolume += klines[i].Volume
	}
	for i := midPoint; i < len(klines); i++ {
		secondHalfVolume += klines[i].Volume
	}

	firstHalfAvg := firstHalfVolume / float64(midPoint)
	secondHalfAvg := secondHalfVolume / float64(len(klines)-midPoint)

	var volumeTrend analyticsTypes.TrendDirection
	trendThreshold := 1.2 // 20% increase/decrease

	if firstHalfAvg > 0 && secondHalfAvg/firstHalfAvg > trendThreshold {
		volumeTrend = analyticsTypes.TrendBullish // Increasing volume
	} else if secondHalfAvg > 0 && firstHalfAvg/secondHalfAvg > trendThreshold {
		volumeTrend = analyticsTypes.TrendBearish // Decreasing volume
	} else {
		volumeTrend = analyticsTypes.TrendNeutral
	}

	return &analyticsTypes.VolumeAnalysis{
		CurrentVolume: numerical.NewFromFloat(currentVolume),
		AverageVolume: numerical.NewFromFloat(avgVolume),
		VolumeRatio:   numerical.NewFromFloat(volumeRatio),
		IsVolumeSpike: isSpike,
		VolumeTrend:   volumeTrend,
	}, nil
}
