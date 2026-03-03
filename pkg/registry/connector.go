package registry

import (
	"fmt"
	"sync"
	"time"

	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/connector/spot"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

type connectorState struct {
	connector     connector.Connector
	connectorType connector.MarketType
	ready         bool
	readyAt       time.Time
}

type connectorRegistry struct {
	connectors map[connector.ExchangeName]*connectorState
	mu         sync.RWMutex
}

func NewConnectorRegistry() registry.ConnectorRegistry {
	return &connectorRegistry{
		connectors: make(map[connector.ExchangeName]*connectorState),
	}
}

// ===== Registration =====

func (cr *connectorRegistry) RegisterSpot(name connector.ExchangeName, conn spot.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.connectors[name] = &connectorState{
		connector:     conn,
		connectorType: connector.MarketTypeSpot,
		ready:         false,
	}
}

func (cr *connectorRegistry) RegisterPerp(name connector.ExchangeName, conn perp.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.connectors[name] = &connectorState{
		connector:     conn,
		connectorType: connector.MarketTypePerp,
		ready:         false,
	}
}

func (cr *connectorRegistry) RegisterPrediction(name connector.ExchangeName, conn predictionconnector.Connector) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.connectors[name] = &connectorState{
		connector:     conn,
		connectorType: connector.MarketTypePrediction,
		ready:         false,
	}
}

// ===== Direct Getters =====

func (cr *connectorRegistry) Connector(name connector.ExchangeName) (connector.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists {
		return nil, false
	}
	return state.connector, true
}

func (cr *connectorRegistry) Spot(name connector.ExchangeName) (spot.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists || state.connectorType != connector.MarketTypeSpot {
		return nil, false
	}

	if spotConn, ok := state.connector.(spot.Connector); ok {
		return spotConn, true
	}
	return nil, false
}

func (cr *connectorRegistry) Perp(name connector.ExchangeName) (perp.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists || state.connectorType != connector.MarketTypePerp {
		return nil, false
	}

	if perpConn, ok := state.connector.(perp.Connector); ok {
		return perpConn, true
	}
	return nil, false
}

func (cr *connectorRegistry) Prediction(name connector.ExchangeName) (predictionconnector.Connector, bool) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	if !exists || state.connectorType != connector.MarketTypePrediction {
		return nil, false
	}

	if predConn, ok := state.connector.(predictionconnector.Connector); ok {
		return predConn, true
	}
	return nil, false
}

// ===== WebSocket Getters =====

func (cr *connectorRegistry) SpotWebSocket(name connector.ExchangeName) (spot.WebSocketConnector, bool) {
	conn, ok := cr.Spot(name)
	if !ok {
		return nil, false
	}

	if ws, ok := conn.(spot.WebSocketConnector); ok {
		return ws, true
	}
	return nil, false
}

func (cr *connectorRegistry) PerpWebSocket(name connector.ExchangeName) (perp.WebSocketConnector, bool) {
	conn, ok := cr.Perp(name)
	if !ok {
		return nil, false
	}

	if ws, ok := conn.(perp.WebSocketConnector); ok {
		return ws, true
	}
	return nil, false
}

func (cr *connectorRegistry) PredictionWebSocket(name connector.ExchangeName) (predictionconnector.WebSocketConnector, bool) {
	conn, ok := cr.Prediction(name)
	if !ok {
		return nil, false
	}

	if ws, ok := conn.(predictionconnector.WebSocketConnector); ok {
		return ws, true
	}
	return nil, false
}

// ===== Filter-based Queries =====

func (cr *connectorRegistry) Filter(opts registry.FilterOptions) []connector.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var results []connector.Connector
	for _, state := range cr.connectors {
		if cr.matchesFilter(state, opts, connector.MarketType("")) {
			results = append(results, state.connector)
		}
	}
	return results
}

func (cr *connectorRegistry) FilterSpot(opts registry.FilterOptions) []spot.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var results []spot.Connector
	for _, state := range cr.connectors {
		if state.connectorType != connector.MarketTypeSpot {
			continue
		}
		if cr.matchesFilter(state, opts, connector.MarketTypeSpot) {
			if spotConn, ok := state.connector.(spot.Connector); ok {
				results = append(results, spotConn)
			}
		}
	}
	return results
}

func (cr *connectorRegistry) FilterPerp(opts registry.FilterOptions) []perp.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var results []perp.Connector
	for _, state := range cr.connectors {
		if state.connectorType != connector.MarketTypePerp {
			continue
		}
		if cr.matchesFilter(state, opts, connector.MarketTypePerp) {
			if perpConn, ok := state.connector.(perp.Connector); ok {
				results = append(results, perpConn)
			}
		}
	}
	return results
}

func (cr *connectorRegistry) FilterPrediction(opts registry.FilterOptions) []predictionconnector.Connector {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	var results []predictionconnector.Connector
	for _, state := range cr.connectors {
		if state.connectorType != connector.MarketTypePrediction {
			continue
		}
		if cr.matchesFilter(state, opts, connector.MarketTypePrediction) {
			if predConn, ok := state.connector.(predictionconnector.Connector); ok {
				results = append(results, predConn)
			}
		}
	}
	return results
}

// ===== Filter Helper =====

func (cr *connectorRegistry) matchesFilter(state *connectorState, opts registry.FilterOptions, marketType connector.MarketType) bool {
	// Check ready state
	if opts.IsReadyOnly() && !state.ready {
		return false
	}

	// Check websocket capability
	if opts.IsWebSocketOnly() {
		switch marketType {
		case connector.MarketTypeSpot:
			_, ok := state.connector.(spot.WebSocketConnector)
			return ok
		case connector.MarketTypePerp:
			_, ok := state.connector.(perp.WebSocketConnector)
			return ok
		case connector.MarketTypePrediction:
			_, ok := state.connector.(predictionconnector.WebSocketConnector)
			return ok
		default:
			// For generic queries, check any websocket interface
			_, spotWS := state.connector.(spot.WebSocketConnector)
			_, perpWS := state.connector.(perp.WebSocketConnector)
			return spotWS || perpWS
		}
	}

	return true
}

// ===== Ready State Management =====

func (cr *connectorRegistry) MarkReady(name connector.ExchangeName) error {
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

func (cr *connectorRegistry) IsReady(name connector.ExchangeName) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	state, exists := cr.connectors[name]
	return exists && state.ready
}
