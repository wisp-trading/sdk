package activity

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/shopspring/decimal"
)

// Trades manages trade execution history globally (trades contain their exchange info)
type Trades interface {
	// Add new trades (from order execution or exchange sync)
	AddTrade(trade connector.Trade)
	AddTrades(trades []connector.Trade)

	// Query trades globally (trades contain exchange in their struct)
	GetAllTrades() []connector.Trade
	GetTradesByExchange(exchange connector.ExchangeName) []connector.Trade
	GetTradesByAsset(asset portfolio.Asset) []connector.Trade
	GetTradesByExchangeAndAsset(exchange connector.ExchangeName, asset portfolio.Asset) []connector.Trade
	GetTradesSince(since time.Time) []connector.Trade

	// Query by ID
	GetTradeByID(tradeID string) *connector.Trade
	TradeExists(tradeID string) bool

	// Analytics helpers
	GetTradeCount() int
	GetTotalVolume(asset portfolio.Asset) decimal.Decimal

	// Clear for simulation restart
	Clear()
}
