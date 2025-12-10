package activity

import (
	"context"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Trades provides read-only access to trade history
type Trades interface {
	// Query trades
	GetAllTrades(ctx context.Context) []connector.Trade
	GetTradesByExchange(ctx context.Context, exchange connector.ExchangeName) []connector.Trade
	GetTradesByAsset(ctx context.Context, asset portfolio.Asset) []connector.Trade
	GetTradesSince(ctx context.Context, since time.Time) []connector.Trade
	GetTradeByID(ctx context.Context, tradeID string) *connector.Trade

	// Stats
	GetTradeCount(ctx context.Context) int
	GetTotalVolume(ctx context.Context, asset portfolio.Asset) numerical.Decimal
}
