package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// MarketDataReader provides read-only market data access
type MarketDataReader interface {
	FetchPrice(symbol string) (*Price, error)
	FetchKlines(symbol, interval string, limit int) ([]Kline, error)
	FetchOrderBook(symbol portfolio.Asset, depth int) (*OrderBook, error)
	FetchRecentTrades(symbol string, limit int) ([]Trade, error)
}

// OrderExecutor handles order placement and management
type OrderExecutor interface {
	PlaceLimitOrder(symbol string, side OrderSide, quantity, price numerical.Decimal) (*OrderResponse, error)
	PlaceMarketOrder(symbol string, side OrderSide, quantity numerical.Decimal) (*OrderResponse, error)
	CancelOrder(symbol, orderID string) (*CancelResponse, error)
	GetOpenOrders() ([]Order, error)
	GetOrderStatus(orderID string) (*Order, error)
}

// AccountReader provides account information
type AccountReader interface {
	GetAccountBalance() (*AccountBalance, error)
	GetTradingHistory(symbol string, limit int) ([]Trade, error)
}

// WebSocketCapable provides WebSocket lifecycle management
type WebSocketCapable interface {
	StartWebSocket() error
	StopWebSocket() error
	IsWebSocketConnected() bool
	ErrorChannel() <-chan error
}
