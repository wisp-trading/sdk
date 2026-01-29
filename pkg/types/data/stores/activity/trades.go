package activity

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Trades manages trade execution history globally (trades contain their exchange info)
type Trades interface {
	// Add new trades
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
	GetTotalVolume(asset portfolio.Asset) numerical.Decimal

	// Clear for simulation restart
	Clear()
}
