package analytics

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Market provides market data access and analysis.
type Market interface {
	// Asset data retrieval
	GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Asset
	FundingRates(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]perp.FundingRate
	FundingRate(ctx context.Context, asset portfolio.Asset, exchange connector.ExchangeName) (*perp.FundingRate, error)

	// Price data
	Price(ctx context.Context, asset portfolio.Asset, opts ...MarketOptions) (numerical.Decimal, error)
	Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal

	// Order book
	OrderBook(ctx context.Context, asset portfolio.Asset, opts ...MarketOptions) (*connector.OrderBook, error)

	// Liquidity
	GetTradableQuantity(ctx context.Context, asset portfolio.Asset, opts ...LiquidityOptions) numerical.Decimal

	// Arbitrage
	FindArbitrage(ctx context.Context, asset portfolio.Asset, minSpreadBps numerical.Decimal) []ArbitrageOpportunity
}

// MarketOptions configures market data queries
type MarketOptions struct {
	Exchange       connector.ExchangeName // Optional: defaults to first available exchange
	InstrumentType connector.Instrument   // Optional: defaults to perpetual
}

// LiquidityOptions configures how liquidity is calculated
type LiquidityOptions struct {
	MaxOrderSizeUSD    numerical.Decimal // Maximum order size in USD
	LiquidityDepthPct  numerical.Decimal // Percentage of order book depth to use (e.g., 0.1 = 10%)
	MinLiquidityLevels int               // Minimum number of price levels to check
}

// ArbitrageOpportunity represents a price discrepancy across exchanges
type ArbitrageOpportunity struct {
	Asset         portfolio.Asset
	BuyExchange   connector.ExchangeName
	SellExchange  connector.ExchangeName
	BuyPrice      numerical.Decimal
	SellPrice     numerical.Decimal
	SpreadBps     numerical.Decimal // Spread in basis points
	SpreadPercent numerical.Decimal // Spread as percentage
}
