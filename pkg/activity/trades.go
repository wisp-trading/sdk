package activity

import (
	"context"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	storeActivity "github.com/wisp-trading/sdk/pkg/types/data/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// trades wraps the internal trade store with read-only access
type trades struct {
	store storeActivity.Trades
}

// NewTrades creates a new read-only trades accessor
func NewTrades(store storeActivity.Trades) wispActivity.Trades {
	return &trades{store: store}
}

// GetAllTrades retrieves all trades
func (t *trades) GetAllTrades(ctx context.Context) []connector.Trade {
	return t.store.GetAllTrades()
}

// GetTradesByExchange retrieves trades for a specific exchange
func (t *trades) GetTradesByExchange(ctx context.Context, exchange connector.ExchangeName) []connector.Trade {
	return t.store.GetTradesByExchange(exchange)
}

// GetTradesByAsset retrieves trades for a specific asset
func (t *trades) GetTradesByAsset(ctx context.Context, asset portfolio.Pair) []connector.Trade {
	return t.store.GetTradesByPair(asset)
}

// GetTradesSince retrieves trades since a specific time
func (t *trades) GetTradesSince(ctx context.Context, since time.Time) []connector.Trade {
	return t.store.GetTradesSince(since)
}

// GetTradeByID retrieves a trade by ID
func (t *trades) GetTradeByID(ctx context.Context, tradeID string) *connector.Trade {
	return t.store.GetTradeByID(tradeID)
}

// GetTradeCount returns the total number of trades
func (t *trades) GetTradeCount(ctx context.Context) int {
	return t.store.GetTradeCount()
}

// GetTotalVolume calculates total volume for a specific asset
func (t *trades) GetTotalVolume(ctx context.Context, asset portfolio.Pair) numerical.Decimal {
	return t.store.GetTotalVolume(asset)
}

var _ wispActivity.Trades = (*trades)(nil)
