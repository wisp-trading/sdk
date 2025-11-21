package registry

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
)

type ConnectorRegistry interface {
	GetConnector(name connector.ExchangeName) (connector.Connector, bool)
	RegisterConnector(name connector.ExchangeName, conn connector.Connector)
	RegisterAllConnectors(connectors []connector.Connector)
	GetAvailableConnectors() []connector.Connector
	GetWebSocketConnectors() []connector.WebSocketConnector
	GetTradingWebSocketConnectors() []connector.WebSocketConnector

	// Readiness tracking
	MarkConnectorReady(name connector.ExchangeName) error
	IsConnectorReady(name connector.ExchangeName) bool
	GetReadyConnectors() []connector.Connector

	// GetDataTimeRange returns the time range of historical data if available
	GetDataTimeRange() (start, end time.Time, err error)
}
