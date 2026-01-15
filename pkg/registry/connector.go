package registry

import (
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector/spot"
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

// ===== Spot Connector Methods =====

func (cr *connectorRegistry) GetSpotConnector(name connector.ExchangeName) (spot.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists {
		return nil, false
	}

	if spotConn, ok := state.connector.(spot.Connector); ok {
		return spotConn, true
	}

	return nil, false
}

func (cr *connectorRegistry) RegisterSpotConnector(name connector.ExchangeName, conn spot.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.connectors[name] = &connectorState{
		connector: conn,
		ready:     false,
	}
}

func (cr *connectorRegistry) GetSpotConnectors() []spot.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var spotConnectors []spot.Connector
	for _, state := range cr.connectors {
		if spotConn, ok := state.connector.(spot.Connector); ok {
			spotConnectors = append(spotConnectors, spotConn)
		}
	}

	return spotConnectors
}

func (cr *connectorRegistry) GetReadySpotConnectors() []spot.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var spotConnectors []spot.Connector
	for _, state := range cr.connectors {
		if state.ready {
			if spotConn, ok := state.connector.(spot.Connector); ok {
				spotConnectors = append(spotConnectors, spotConn)
			}
		}
	}

	return spotConnectors
}

func (cr *connectorRegistry) GetSpotWebSocketConnectors() []spot.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []spot.WebSocketConnector
	for _, state := range cr.connectors {
		if wsConn, ok := state.connector.(spot.WebSocketConnector); ok {
			wsConnectors = append(wsConnectors, wsConn)
		}
	}

	return wsConnectors
}

func (cr *connectorRegistry) GetReadySpotWebSocketConnectors() []spot.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []spot.WebSocketConnector
	for _, state := range cr.connectors {
		if state.ready {
			if wsConn, ok := state.connector.(spot.WebSocketConnector); ok {
				wsConnectors = append(wsConnectors, wsConn)
			}
		}
	}

	return wsConnectors
}

// ===== Perpetual Connector Methods =====

func (cr *connectorRegistry) GetPerpConnector(name connector.ExchangeName) (perp.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists {
		return nil, false
	}

	if perpConn, ok := state.connector.(perp.Connector); ok {
		return perpConn, true
	}

	return nil, false
}

func (cr *connectorRegistry) RegisterPerpConnector(name connector.ExchangeName, conn perp.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.connectors[name] = &connectorState{
		connector: conn,
		ready:     false,
	}
}

func (cr *connectorRegistry) GetPerpConnectors() []perp.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var perpConnectors []perp.Connector
	for _, state := range cr.connectors {
		if perpConn, ok := state.connector.(perp.Connector); ok {
			perpConnectors = append(perpConnectors, perpConn)
		}
	}

	return perpConnectors
}

func (cr *connectorRegistry) GetReadyPerpConnectors() []perp.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var perpConnectors []perp.Connector
	for _, state := range cr.connectors {
		if state.ready {
			if perpConn, ok := state.connector.(perp.Connector); ok {
				perpConnectors = append(perpConnectors, perpConn)
			}
		}
	}

	return perpConnectors
}

func (cr *connectorRegistry) GetPerpWebSocketConnectors() []perp.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []perp.WebSocketConnector
	for _, state := range cr.connectors {
		if wsConn, ok := state.connector.(perp.WebSocketConnector); ok {
			wsConnectors = append(wsConnectors, wsConn)
		}
	}

	return wsConnectors
}

func (cr *connectorRegistry) GetReadyPerpWebSocketConnectors() []perp.WebSocketConnector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var wsConnectors []perp.WebSocketConnector
	for _, state := range cr.connectors {
		if state.ready {
			if wsConn, ok := state.connector.(perp.WebSocketConnector); ok {
				wsConnectors = append(wsConnectors, wsConn)
			}
		}
	}

	return wsConnectors
}

// ===== Generic Base Connector Methods =====

func (cr *connectorRegistry) GetConnector(name connector.ExchangeName) (connector.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists {
		return nil, false
	}

	// Try spot first
	if spotConn, ok := state.connector.(spot.Connector); ok {
		return spotConn, true
	}

	// Try perp
	if perpConn, ok := state.connector.(perp.Connector); ok {
		return perpConn, true
	}

	return nil, false
}

func (cr *connectorRegistry) GetAllBaseConnectors() []connector.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var baseConnectors []connector.Connector
	for _, state := range cr.connectors {
		if spotConn, ok := state.connector.(spot.Connector); ok {
			baseConnectors = append(baseConnectors, spotConn)
		} else if perpConn, ok := state.connector.(perp.Connector); ok {
			baseConnectors = append(baseConnectors, perpConn)
		}
	}

	return baseConnectors
}

func (cr *connectorRegistry) GetAllReadyConnectors() []connector.Connector {
	var readyConnectors []connector.Connector

	// Convert perp connectors to Connector
	perpConnectors := cr.GetReadyPerpConnectors()
	for _, conn := range perpConnectors {
		readyConnectors = append(readyConnectors, conn)
	}

	// Convert spot connectors to Connector
	spotConnectors := cr.GetReadySpotConnectors()
	for _, conn := range spotConnectors {
		readyConnectors = append(readyConnectors, conn)
	}

	return readyConnectors
}

// ===== Ready State Management =====

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

func (cr *connectorRegistry) IsConnectorReady(name connector.ExchangeName) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	return exists && state.ready
}
