package common

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// MarketDataReader provides read-only market data access
type MarketDataReader interface {
	FetchPrice(symbol string) (*connector.Price, error)
	FetchKlines(symbol, interval string, limit int) ([]connector.Kline, error)
	FetchOrderBook(symbol portfolio.Asset, depth int) (*connector.OrderBook, error)
	FetchRecentTrades(symbol string, limit int) ([]connector.Trade, error)
}

// OrderExecutor handles order placement and management
type OrderExecutor interface {
	PlaceLimitOrder(symbol string, side connector.OrderSide, quantity, price numerical.Decimal) (*connector.OrderResponse, error)
	PlaceMarketOrder(symbol string, side connector.OrderSide, quantity numerical.Decimal) (*connector.OrderResponse, error)
	CancelOrder(symbol, orderID string) (*connector.CancelResponse, error)
	GetOpenOrders() ([]connector.Order, error)
	GetOrderStatus(orderID string) (*connector.Order, error)
}

// AccountReader provides account information
type AccountReader interface {
	GetAccountBalance() (*connector.AccountBalance, error)
	GetTradingHistory(symbol string, limit int) ([]connector.Trade, error)
}

// WebSocketCapable provides WebSocket lifecycle management
type WebSocketCapable interface {
	StartWebSocket() error
	StopWebSocket() error
	IsWebSocketConnected() bool
	ErrorChannel() <-chan error
}
