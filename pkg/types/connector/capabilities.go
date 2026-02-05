package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// MarketDataReader provides read-only market data access
type MarketDataReader interface {
	FetchPrice(pair portfolio.Pair) (*Price, error)
	FetchKlines(pair portfolio.Pair, interval string, limit int) ([]Kline, error)
	FetchOrderBook(pair portfolio.Pair, depth int) (*OrderBook, error)
	FetchRecentTrades(pair portfolio.Pair, limit int) ([]Trade, error)
}

// OrderExecutor handles order placement and management
type OrderExecutor interface {
	PlaceLimitOrder(pair portfolio.Pair, side OrderSide, quantity, price numerical.Decimal) (*OrderResponse, error)
	PlaceMarketOrder(pair portfolio.Pair, side OrderSide, quantity numerical.Decimal) (*OrderResponse, error)
	CancelOrder(orderID string, pair ...portfolio.Pair) (*CancelResponse, error)
	GetOpenOrders(pair ...portfolio.Pair) ([]Order, error)
	GetOrderStatus(orderID string, pair ...portfolio.Pair) (*Order, error)
}

// AccountReader provides account information
type AccountReader interface {
	GetBalances() ([]AssetBalance, error)
	GetBalance(asset portfolio.Asset) (*AssetBalance, error)
	GetTradingHistory(pair portfolio.Pair, limit int) ([]Trade, error)
}

// WebSocketCapable provides WebSocket lifecycle management
type WebSocketCapable interface {
	StartWebSocket() error
	StopWebSocket() error
	IsWebSocketConnected() bool
	ErrorChannel() <-chan error
}
