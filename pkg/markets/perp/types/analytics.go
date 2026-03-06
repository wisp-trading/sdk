package types

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PerpMarket provides perpetual-specific market data access
type PerpMarket interface {
	// Price returns the current price for an asset on perp exchanges
	Price(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error)

	// Prices returns prices across all perp exchanges
	Prices(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal

	// OrderBook returns the order book for an asset on a perp exchange
	OrderBook(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*connector.OrderBook, error)

	// GetKlines returns historical kline/candlestick data for perp markets
	GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline

	// GetTradableQuantity calculates available liquidity for perp trading
	GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analytics.LiquidityOptions) numerical.Decimal

	// FundingRate returns the funding rate for an asset on a specific perp exchange
	FundingRate(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*perp.FundingRate, error)

	// FundingRates returns funding rates across all perp exchanges
	FundingRates(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]perp.FundingRate

	// GetAllAssetsWithFundingRates returns all assets that have funding rate data
	GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Pair
}
