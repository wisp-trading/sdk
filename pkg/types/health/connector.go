package health

import "github.com/backtesting-org/kronos-sdk/pkg/types/connector"

// ConnectorErrorStore tracks WebSocket connector-level errors and connection state
type ConnectorErrorStore interface {
	RecordConnectorError(name connector.ExchangeName, err error)
	UpdateConnectionState(name connector.ExchangeName, state ConnectionState)
	GetConnectorState(name connector.ExchangeName) (ConnectionState, bool)
	GetConnectorError(name connector.ExchangeName) error
	CountTrackedConnectors() int
	GetUnhealthyConnectors() []connector.ExchangeName
}

// ConnectionState represents the connection status
type ConnectionState string

const (
	StateConnected    ConnectionState = "connected"
	StateDisconnected ConnectionState = "disconnected"
	StateConnecting   ConnectionState = "connecting"
	StateDegraded     ConnectionState = "degraded"
)
