package analytics

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// VolumeAnalysis detects volume patterns and spikes from the provided klines.
func (s *analytics) VolumeAnalysis(klines []connector.Kline) (*analyticsTypes.VolumeAnalysis, error) {
	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for volume analysis")
	}

	currentVolume := klines[len(klines)-1].Volume

	totalVolume := 0.0
	for i := 0; i < len(klines)-1; i++ {
		totalVolume += klines[i].Volume
	}
	avgVolume := totalVolume / float64(len(klines)-1)

	var volumeRatio float64
	if avgVolume > 0 {
		volumeRatio = currentVolume / avgVolume
	}

	isSpike := volumeRatio > 2.0

	midPoint := len(klines) / 2
	var firstHalf, secondHalf float64
	for i := 0; i < midPoint; i++ {
		firstHalf += klines[i].Volume
	}
	for i := midPoint; i < len(klines); i++ {
		secondHalf += klines[i].Volume
	}

	firstHalfAvg := firstHalf / float64(midPoint)
	secondHalfAvg := secondHalf / float64(len(klines)-midPoint)

	const trendThreshold = 1.2
	var volumeTrend analyticsTypes.TrendDirection
	if firstHalfAvg > 0 && secondHalfAvg/firstHalfAvg > trendThreshold {
		volumeTrend = analyticsTypes.TrendBullish
	} else if secondHalfAvg > 0 && firstHalfAvg/secondHalfAvg > trendThreshold {
		volumeTrend = analyticsTypes.TrendBearish
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
