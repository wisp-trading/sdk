package activity

import (
	"context"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Trades provides read-only access to trade history
type Trades interface {
	// Query trades
	GetAllTrades(ctx context.Context) []connector.Trade
	GetTradesByExchange(ctx context.Context, exchange connector.ExchangeName) []connector.Trade
	GetTradesByAsset(ctx context.Context, asset portfolio.Pair) []connector.Trade
	GetTradesSince(ctx context.Context, since time.Time) []connector.Trade
	GetTradeByID(ctx context.Context, tradeID string) *connector.Trade

	// Stats
	GetTradeCount(ctx context.Context) int
	GetTotalVolume(ctx context.Context, asset portfolio.Pair) numerical.Decimal
}
