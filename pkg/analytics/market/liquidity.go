package market

import (
	"context"
	"time"

	"github.com/wisp-trading/sdk/pkg/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// DefaultLiquidityOptions returns sensible defaults
func DefaultLiquidityOptions() analytics.LiquidityOptions {
	return analytics.LiquidityOptions{
		MaxOrderSizeUSD:    numerical.NewFromInt(10000), // $10k default
		LiquidityDepthPct:  numerical.NewFromFloat(0.1), // Use 10% of available liquidity
		MinLiquidityLevels: 5,                           // Check top 5 levels
	}
}

// getTradableQuantity calculates the maximum tradable quantity based on order book liquidity
// Returns the quantity in base currency that can be safely traded
// This is a shared helper used by both spot and perp market services
func getTradableQuantity(ctx context.Context, orderBook *connector.OrderBook, opts analytics.LiquidityOptions) numerical.Decimal {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("GetTradableQuantity", time.Since(start))
		}
	}()

	if orderBook == nil {
		return numerical.Zero()
	}

	// Check if we have sufficient depth
	if len(orderBook.Bids) < opts.MinLiquidityLevels || len(orderBook.Asks) < opts.MinLiquidityLevels {
		return numerical.Zero()
	}

	// Calculate available liquidity on both sides
	bidLiquidity := calculateSideLiquidity(orderBook.Bids, opts.MinLiquidityLevels)
	askLiquidity := calculateSideLiquidity(orderBook.Asks, opts.MinLiquidityLevels)

	// Use the smaller of the two (bottleneck)
	availableLiquidity := bidLiquidity
	if askLiquidity.LessThan(bidLiquidity) {
		availableLiquidity = askLiquidity
	}

	// Apply liquidity depth percentage (only use a fraction of available liquidity)
	usableLiquidity := availableLiquidity.Mul(opts.LiquidityDepthPct)

	// Get mid price for USD conversion
	midPrice := orderBook.Bids[0].Price.Add(orderBook.Asks[0].Price).Div(numerical.NewFromInt(2))

	// Calculate max quantity in base currency based on USD limit
	maxQuantityUSD := opts.MaxOrderSizeUSD.Div(midPrice)

	// Return the smaller of usable liquidity or max order size
	if usableLiquidity.LessThan(maxQuantityUSD) {
		return usableLiquidity
	}
	return maxQuantityUSD
}

// calculateSideLiquidity sums up the quantity available in the order book side
func calculateSideLiquidity(levels []connector.PriceLevel, maxLevels int) numerical.Decimal {
	total := numerical.Zero()
	levelsToUse := maxLevels
	if len(levels) < levelsToUse {
		levelsToUse = len(levels)
	}

	for i := 0; i < levelsToUse; i++ {
		total = total.Add(levels[i].Quantity)
	}

	return total
}
