package types

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// SpotMarket provides spot-specific market data access
type SpotMarket interface {
	// Price returns the current price for an asset on spot exchanges
	Price(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error)

	// Prices returns prices across all spot exchanges
	Prices(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal

	// OrderBook returns the order book for an asset on a spot exchange
	OrderBook(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*connector.OrderBook, error)

	// GetKlines returns historical kline/candlestick data for spot markets
	GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline

	// GetTradableQuantity calculates available liquidity for spot trading
	GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analytics.LiquidityOptions) numerical.Decimal
}
