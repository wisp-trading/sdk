package analytics

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Market provides market data access and analysis across all market types.
// Common methods work across both spot and perp markets.
// Use Spot() or Perp() to access market-type-specific functionality.
type Market interface {
	// Spot returns spot-specific market service
	Spot() SpotMarket

	// Perp returns perpetual-specific market service
	Perp() PerpMarket

	// Price returns the current price for an asset from specified exchange or first available.
	// Works across both spot and perp exchanges.
	Price(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (numerical.Decimal, error)

	// Prices returns prices for an asset across all spot and perp exchanges.
	Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal

	// Klines returns historical kline data for an asset on the specified exchange.
	// Automatically searches all registered market stores (spot, perp, etc.) to find the exchange.
	// The user doesn't need to know which market type the exchange belongs to.
	Klines(asset portfolio.Asset, exchange connector.ExchangeName, interval string, limit int) []connector.Kline

	// OrderBook returns the order book for an asset on the specified exchange.
	// Automatically searches all registered market stores to find which one has this exchange.
	// The user doesn't need to know which market type the exchange belongs to.
	OrderBook(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (*connector.OrderBook, error)

	// FindArbitrage finds arbitrage opportunities across all exchanges (spot and perp).
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

// SpotMarket provides spot-specific market data access
type SpotMarket interface {
	// Price returns the current price for an asset on spot exchanges
	Price(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (numerical.Decimal, error)

	// Prices returns prices across all spot exchanges
	Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal

	// OrderBook returns the order book for an asset on a spot exchange
	OrderBook(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (*connector.OrderBook, error)

	// GetKlines returns historical kline/candlestick data for spot markets
	GetKlines(asset portfolio.Asset, exchange connector.ExchangeName, interval string, limit int) []connector.Kline

	// GetTradableQuantity calculates available liquidity for spot trading
	GetTradableQuantity(ctx context.Context, asset portfolio.Asset, opts ...LiquidityOptions) numerical.Decimal
}

// PerpMarket provides perpetual-specific market data access
type PerpMarket interface {
	// Price returns the current price for an asset on perp exchanges
	Price(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (numerical.Decimal, error)

	// Prices returns prices across all perp exchanges
	Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal

	// OrderBook returns the order book for an asset on a perp exchange
	OrderBook(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (*connector.OrderBook, error)

	// GetKlines returns historical kline/candlestick data for perp markets
	GetKlines(asset portfolio.Asset, exchange connector.ExchangeName, interval string, limit int) []connector.Kline

	// GetTradableQuantity calculates available liquidity for perp trading
	GetTradableQuantity(ctx context.Context, asset portfolio.Asset, opts ...LiquidityOptions) numerical.Decimal

	// FundingRate returns the funding rate for an asset on a specific perp exchange
	FundingRate(ctx context.Context, asset portfolio.Asset, exchange connector.ExchangeName) (*perp.FundingRate, error)

	// FundingRates returns funding rates across all perp exchanges
	FundingRates(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]perp.FundingRate

	// GetAllAssetsWithFundingRates returns all assets that have funding rate data
	GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Asset
}
