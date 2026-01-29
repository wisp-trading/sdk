package registry

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/connector/spot"
)

type ConnectorRegistry interface {
	// Spot connector operations
	GetSpotConnector(name connector.ExchangeName) (spot.Connector, bool)
	RegisterSpotConnector(name connector.ExchangeName, conn spot.Connector)
	GetSpotConnectors() []spot.Connector
	GetReadySpotConnectors() []spot.Connector
	GetSpotWebSocketConnectors() []spot.WebSocketConnector
	GetReadySpotWebSocketConnectors() []spot.WebSocketConnector

	// Perpetual connector operations
	GetPerpConnector(name connector.ExchangeName) (perp.Connector, bool)
	RegisterPerpConnector(name connector.ExchangeName, conn perp.Connector)
	GetPerpConnectors() []perp.Connector
	GetReadyPerpConnectors() []perp.Connector
	GetPerpWebSocketConnectors() []perp.WebSocketConnector
	GetReadyPerpWebSocketConnectors() []perp.WebSocketConnector

	// Generic access (returns base interface)
	GetConnector(name connector.ExchangeName) (connector.Connector, bool)
	GetAllBaseConnectors() []connector.Connector
	GetAllReadyConnectors() []connector.Connector

	// Ready state management (works for all connector types)
	MarkConnectorReady(name connector.ExchangeName) error
	IsConnectorReady(name connector.ExchangeName) bool
}
