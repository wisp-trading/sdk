package connector

import (
	"context"
	"kronos/sdk/pkg/types/portfolio"
)

// WebSocketConnector extends the base Connector with real-time capabilities
type WebSocketConnector interface {
	Connector

	// WebSocket lifecycle management
	StartWebSocket(ctx context.Context) error
	StopWebSocket() error
	IsWebSocketConnected() bool

	// Subscription management
	SubscribeOrderBook(asset portfolio.Asset, instrumentType Instrument) error
	SubscribeTrades(asset portfolio.Asset, instrumentType Instrument) error
	SubscribePositions(asset portfolio.Asset, instrumentType Instrument) error
	SubscribeAccountBalance() error
	SubscribeKlines(asset portfolio.Asset, interval string) error

	UnsubscribeKlines(asset portfolio.Asset, interval string) error
	UnsubscribeTrades(asset portfolio.Asset, instrumentType Instrument) error
	UnsubscribeOrderBook(asset portfolio.Asset, instrumentType Instrument) error
	UnsubscribePositions(asset portfolio.Asset, instrumentType Instrument) error
	UnsubscribeAccountBalance() error

	// Data channels - all required for RealtimeIngestor
	OrderBookUpdates() <-chan OrderBook
	TradeUpdates() <-chan Trade
	PositionUpdates() <-chan Position
	AccountBalanceUpdates() <-chan AccountBalance
	KlineUpdates() <-chan Kline
	ErrorChannel() <-chan error
}
