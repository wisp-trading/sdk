package registry

import (
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/connector/spot"
)

type ConnectorRegistry interface {
	// Registration
	RegisterSpot(name connector.ExchangeName, conn spot.Connector)
	RegisterPerp(name connector.ExchangeName, conn perp.Connector)
	RegisterPrediction(name connector.ExchangeName, conn predictionconnector.Connector)

	Connector(name connector.ExchangeName) (connector.Connector, bool)

	// ConnectorType returns the market type of a registered connector.
	ConnectorType(name connector.ExchangeName) (connector.MarketType, bool)

	// Direct getters (type-specific)
	Spot(name connector.ExchangeName) (spot.Connector, bool)
	Perp(name connector.ExchangeName) (perp.Connector, bool)
	Prediction(name connector.ExchangeName) (predictionconnector.Connector, bool)

	// Filter-based queries
	Filter(opts FilterOptions) []connector.Connector

	// Typed filter helpers
	FilterSpot(opts FilterOptions) []spot.Connector
	FilterPerp(opts FilterOptions) []perp.Connector
	FilterPrediction(opts FilterOptions) []predictionconnector.Connector

	// WebSocket helpers
	SpotWebSocket(name connector.ExchangeName) (spot.WebSocketConnector, bool)
	PerpWebSocket(name connector.ExchangeName) (perp.WebSocketConnector, bool)
	PredictionWebSocket(name connector.ExchangeName) (predictionconnector.WebSocketConnector, bool)

	// Ready state management
	MarkReady(name connector.ExchangeName) error
	IsReady(name connector.ExchangeName) bool
}

type FilterOptions struct {
	readyOnly     bool
	webSocketOnly bool
}

// Builder for FilterOptions
type FilterOptionsBuilder struct {
	opts FilterOptions
}

// NewFilter creates a new filter options builder
func NewFilter() *FilterOptionsBuilder {
	return &FilterOptionsBuilder{}
}

// ReadyOnly filters for ready connectors only
func (b *FilterOptionsBuilder) ReadyOnly() *FilterOptionsBuilder {
	b.opts.readyOnly = true
	return b
}

// WebSocketOnly filters for websocket connectors only
func (b *FilterOptionsBuilder) WebSocketOnly() *FilterOptionsBuilder {
	b.opts.webSocketOnly = true
	return b
}

// Build returns the constructed filter options
func (b *FilterOptionsBuilder) Build() FilterOptions {
	return b.opts
}

// Getters for internal use by registry
func (f FilterOptions) IsReadyOnly() bool {
	return f.readyOnly
}

func (f FilterOptions) IsWebSocketOnly() bool {
	return f.webSocketOnly
}
