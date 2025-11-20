package analytics

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/shopspring/decimal"
)

// Market provides market data access and analysis.
type Market interface {
	// Asset data retrieval
	GetAllAssetsWithFundingRates() []portfolio.Asset
	FundingRates(asset portfolio.Asset) map[connector.ExchangeName]connector.FundingRate
	FundingRate(asset portfolio.Asset, exchange connector.ExchangeName) (*connector.FundingRate, error)

	// Price data
	Price(asset portfolio.Asset, opts ...MarketOptions) (decimal.Decimal, error)
	Prices(asset portfolio.Asset) map[connector.ExchangeName]decimal.Decimal

	// Order book
	OrderBook(asset portfolio.Asset, opts ...MarketOptions) (*connector.OrderBook, error)

	// Liquidity
	GetTradableQuantity(asset portfolio.Asset, opts ...LiquidityOptions) decimal.Decimal

	// Arbitrage
	FindArbitrage(asset portfolio.Asset, minSpreadBps decimal.Decimal) []ArbitrageOpportunity
}

// MarketOptions configures market data queries
type MarketOptions struct {
	Exchange       connector.ExchangeName // Optional: defaults to first available exchange
	InstrumentType connector.Instrument   // Optional: defaults to perpetual
}

// LiquidityOptions configures how liquidity is calculated
type LiquidityOptions struct {
	MaxOrderSizeUSD    decimal.Decimal // Maximum order size in USD
	LiquidityDepthPct  decimal.Decimal // Percentage of order book depth to use (e.g., 0.1 = 10%)
	MinLiquidityLevels int             // Minimum number of price levels to check
}

// ArbitrageOpportunity represents a price discrepancy across exchanges
type ArbitrageOpportunity struct {
	Asset         portfolio.Asset
	BuyExchange   connector.ExchangeName
	SellExchange  connector.ExchangeName
	BuyPrice      decimal.Decimal
	SellPrice     decimal.Decimal
	SpreadBps     decimal.Decimal // Spread in basis points
	SpreadPercent decimal.Decimal // Spread as percentage
}
