package analytics

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// GetPriceChange calculates price statistics from the provided klines.
func (s *analytics) GetPriceChange(klines []connector.Kline) (*analyticsTypes.PriceChange, error) {
	if len(klines) < 2 {
		return nil, fmt.Errorf("insufficient kline data for price change calculation")
	}

	startPrice := klines[0].Open
	endPrice := klines[len(klines)-1].Close

	highPrice := klines[0].High
	lowPrice := klines[0].Low
	for _, k := range klines {
		if k.High > highPrice {
			highPrice = k.High
		}
		if k.Low < lowPrice {
			lowPrice = k.Low
		}
	}

	change := endPrice - startPrice
	changePercent := (change / startPrice) * 100
	priceRange := highPrice - lowPrice
	priceRangePercent := (priceRange / startPrice) * 100

	return &analyticsTypes.PriceChange{
		StartPrice:        numerical.NewFromFloat(startPrice),
		EndPrice:          numerical.NewFromFloat(endPrice),
		Change:            numerical.NewFromFloat(change),
		ChangePercent:     numerical.NewFromFloat(changePercent),
		HighPrice:         numerical.NewFromFloat(highPrice),
		LowPrice:          numerical.NewFromFloat(lowPrice),
		PriceRange:        numerical.NewFromFloat(priceRange),
		PriceRangePercent: numerical.NewFromFloat(priceRangePercent),
	}, nil
}
