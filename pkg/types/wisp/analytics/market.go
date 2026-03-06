package analytics

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// MarketOptions configures market data queries
type MarketOptions struct {
	Exchange       connector.ExchangeName // Optional: defaults to first available exchange
	InstrumentType connector.Instrument   // Optional: defaults to perpetual
}

// LiquidityOptions configures how liquidity is calculated
type LiquidityOptions struct {
	Exchange           connector.ExchangeName // Required: which exchange's order book to use
	MaxOrderSizeUSD    numerical.Decimal      // Maximum order size in USD
	LiquidityDepthPct  numerical.Decimal      // Percentage of order book depth to use (e.g., 0.1 = 10%)
	MinLiquidityLevels int                    // Minimum number of price levels to check
}

// DefaultLiquidityOptions returns sensible defaults.
func DefaultLiquidityOptions() LiquidityOptions {
	return LiquidityOptions{
		MaxOrderSizeUSD:    numerical.NewFromInt(10000),
		LiquidityDepthPct:  numerical.NewFromFloat(0.1),
		MinLiquidityLevels: 5,
	}
}

// CalculateTradableQuantity returns the maximum tradable base quantity given an order book
// and liquidity options. It is a pure function with no store access.
func CalculateTradableQuantity(orderBook *connector.OrderBook, opts LiquidityOptions) numerical.Decimal {
	if orderBook == nil {
		return numerical.Zero()
	}

	if len(orderBook.Bids) < opts.MinLiquidityLevels || len(orderBook.Asks) < opts.MinLiquidityLevels {
		return numerical.Zero()
	}

	bidLiquidity := sumLevels(orderBook.Bids, opts.MinLiquidityLevels)
	askLiquidity := sumLevels(orderBook.Asks, opts.MinLiquidityLevels)

	available := bidLiquidity
	if askLiquidity.LessThan(bidLiquidity) {
		available = askLiquidity
	}

	usable := available.Mul(opts.LiquidityDepthPct)
	midPrice := orderBook.Bids[0].Price.Add(orderBook.Asks[0].Price).Div(numerical.NewFromInt(2))
	maxQty := opts.MaxOrderSizeUSD.Div(midPrice)

	if usable.LessThan(maxQty) {
		return usable
	}
	return maxQty
}

func sumLevels(levels []connector.PriceLevel, max int) numerical.Decimal {
	total := numerical.Zero()
	if len(levels) < max {
		max = len(levels)
	}
	for i := 0; i < max; i++ {
		total = total.Add(levels[i].Quantity)
	}
	return total
}

// ArbitrageOpportunity represents a price discrepancy across exchanges
type ArbitrageOpportunity struct {
	Asset         portfolio.Pair
	BuyExchange   connector.ExchangeName
	SellExchange  connector.ExchangeName
	BuyPrice      numerical.Decimal
	SellPrice     numerical.Decimal
	SpreadBps     numerical.Decimal // Spread in basis points
	SpreadPercent numerical.Decimal // Spread as percentage
}
