package market

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Feature name constants (must match pkg/inference/features/types.go)
const (
	featureMidPrice    = "mid_price"
	featureBidPrice    = "bid_price"
	featureAskPrice    = "ask_price"
	featureLastPrice   = "last_price"
	featureMarkPrice   = "mark_price"
	featureIndexPrice  = "index_price"
	featureFundingRate = "funding_rate"
)

// Extractor computes raw market data features using the market service.
// It extracts basic price data and funding rates from exchange data.
type Extractor struct {
	market analytics.Market
}

// NewExtractor creates a new market data feature extractor.
func NewExtractor(market analytics.Market) *Extractor {
	return &Extractor{
		market: market,
	}
}

// Extract computes market data features and adds them to the feature map.
// Currently supports: mid_price, bid_price, ask_price, last_price,
// mark_price, index_price, funding_rate.
func (e *Extractor) Extract(ctx context.Context, asset portfolio.Pair, featureMap map[string]float64) error {

	// Get last price
	price, err := e.market.Price(ctx, asset)
	if err == nil {
		featureMap[featureLastPrice], _ = price.Float64()
	}

	// Get order book for bid/ask prices
	orderBook, err := e.market.OrderBook(ctx, asset)
	if err == nil && orderBook != nil {
		// Extract bid and ask prices
		if len(orderBook.Bids) > 0 {
			featureMap[featureBidPrice], _ = orderBook.Bids[0].Price.Float64()
		}

		if len(orderBook.Asks) > 0 {
			featureMap[featureAskPrice], _ = orderBook.Asks[0].Price.Float64()
		}

		// Calculate mid price
		if len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
			bidPrice := orderBook.Bids[0].Price
			askPrice := orderBook.Asks[0].Price
			midPrice := bidPrice.Add(askPrice).Div(numerical.NewFromInt(2))
			featureMap[featureMidPrice], _ = midPrice.Float64()
		}
	}

	// Get funding rate data (for perpetual futures only)
	// This provides funding_rate, mark_price, and index_price
	fundingRates := e.market.Perp().FundingRates(ctx, asset)
	if len(fundingRates) > 0 {
		// Use first available funding rate from any perp exchange
		for _, rate := range fundingRates {
			featureMap[featureFundingRate], _ = rate.CurrentRate.Float64()
			featureMap[featureMarkPrice], _ = rate.MarkPrice.Float64()
			featureMap[featureIndexPrice], _ = rate.IndexPrice.Float64()
			break
		}
	}

	return nil
}
