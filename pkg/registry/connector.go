package registry

import (
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

type connectorState struct {
	connector connector.Connector
	ready     bool
	readyAt   time.Time
}

type connectorRegistry struct {
	connectors map[connector.ExchangeName]*connectorState
	mu         sync.RWMutex
}

// NewConnectorRegistry creates a new connector registry
func NewConnectorRegistry() registry.ConnectorRegistry {
	return &connectorRegistry{
		connectors: make(map[connector.ExchangeName]*connectorState),
	}
}

func (cr *connectorRegistry) GetConnector(name connector.ExchangeName) (connector.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists {
		return nil, false
	}
	return state.connector, true
}

func (cr *connectorRegistry) RegisterConnector(name connector.ExchangeName, conn connector.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.connectors[name] = &connectorState{
		connector: conn,
		ready:     false, // Not ready until explicitly marked
	}
}

func (cr *connectorRegistry) RegisterAllConnectors(connectors []connector.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	for _, conn := range connectors {
		name := conn.GetConnectorInfo().Name
		cr.connectors[name] = &connectorState{
			connector: conn,
			ready:     false,
		}
	}
}

func (cr *connectorRegistry) GetAvailableConnectors() []connector.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	connectors := make([]connector.Connector, 0, len(cr.connectors))
	for _, state := range cr.connectors {
		connectors = append(connectors, state.connector)
	}

	return connectors
}

func (cr *connectorRegistry) GetWebSocketConnectors() []connector.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []connector.WebSocketConnector
	for _, state := range cr.connectors {
		if wsConn, ok := state.connector.(connector.WebSocketConnector); ok {
			wsConnectors = append(wsConnectors, wsConn)
		}
	}

	return wsConnectors
}

// MarkConnectorReady marks a connector as ready for use
func (cr *connectorRegistry) MarkConnectorReady(name connector.ExchangeName) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	state, exists := cr.connectors[name]
	if !exists {
		return fmt.Errorf("connector %s not found", name)
	}

	state.ready = true
	state.readyAt = time.Now()
	return nil
}

// IsConnectorReady returns true if the connector is marked as ready
func (cr *connectorRegistry) IsConnectorReady(name connector.ExchangeName) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	return exists && state.ready
}

// GetReadyConnectors returns all connectors that are marked as ready (initialized)
func (cr *connectorRegistry) GetReadyConnectors() []connector.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	connectors := make([]connector.Connector, 0)
	for _, state := range cr.connectors {
		if state.ready {
			connectors = append(connectors, state.connector)
		}
	}

	return connectors
}

// GetReadyWebSocketConnectors returns ready connectors that support WebSocket
func (cr *connectorRegistry) GetReadyWebSocketConnectors() []connector.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []connector.WebSocketConnector
	for _, state := range cr.connectors {
		if state.ready {
			if wsConn, ok := state.connector.(connector.WebSocketConnector); ok {
				wsConnectors = append(wsConnectors, wsConn)
			}
		}
	}

	return wsConnectors
}
