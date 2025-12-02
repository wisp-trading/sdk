package registry

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
)

type ConnectorRegistry interface {
	// Registration phase - connectors are registered but not necessarily initialized
	GetConnector(name connector.ExchangeName) (connector.Connector, bool)
	RegisterConnector(name connector.ExchangeName, conn connector.Connector)
	RegisterAllConnectors(connectors []connector.Connector)
	GetAvailableConnectors() []connector.Connector
	GetWebSocketConnectors() []connector.WebSocketConnector

	// Ready phase - connectors have been initialized and are live
	MarkConnectorReady(name connector.ExchangeName) error
	IsConnectorReady(name connector.ExchangeName) bool
	GetReadyConnectors() []connector.Connector
	GetReadyWebSocketConnectors() []connector.WebSocketConnector
}
