package health

import (
"fmt"
"sync"
"time"

"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
healthTypes "github.com/backtesting-org/kronos-sdk/pkg/types/health"
)

type healthStore struct {
	connectors map[connector.ExchangeName]*healthTypes.ConnectorHealth
	startedAt  time.Time
	mu         sync.RWMutex
}

func NewHealthStore() healthTypes.HealthStore {
	return &healthStore{
		connectors: make(map[connector.ExchangeName]*healthTypes.ConnectorHealth),
		startedAt:  time.Now(),
	}
}

func (h *healthStore) RegisterConnector(name connector.ExchangeName) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.connectors[name]; exists {
		return
	}

	h.connectors[name] = &healthTypes.ConnectorHealth{
		Name:            name,
		State:           healthTypes.StateConnecting,
		DataTypes:       make(map[healthTypes.DataType]*healthTypes.DataTypeHealth),
		LastHealthCheck: time.Now(),
		UptimeSeconds:   0,
		ErrorRate:       0,
	}
}

func (h *healthStore) UpdateConnectionState(name connector.ExchangeName, state healthTypes.ConnectionState) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, exists := h.connectors[name]; exists {
		conn.State = state
		conn.LastHealthCheck = time.Now()
	}
}

func (h *healthStore) RecordDataReceived(name connector.ExchangeName, dataType healthTypes.DataType, source healthTypes.DataSourceType, latency time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &healthTypes.DataTypeHealth{}
	}

	dt := conn.DataTypes[dataType]
	dt.Available = true
	dt.Source = source
	dt.LastReceived = time.Now()
	dt.Latency = latency
	dt.RecordCount++
	dt.ErrorCount = 0
	dt.LastError = nil

	conn.LastHealthCheck = time.Now()
	
	if conn.State != healthTypes.StateConnected {
		conn.State = healthTypes.StateConnected
	}
}

func (h *healthStore) RecordDataError(name connector.ExchangeName, dataType healthTypes.DataType, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &healthTypes.DataTypeHealth{}
	}

	dt := conn.DataTypes[dataType]
	dt.LastError = err
	dt.ErrorCount++

	if dt.ErrorCount >= 3 {
		dt.Available = false
		conn.State = healthTypes.StateDegraded
	}

	conn.LastHealthCheck = time.Now()
}

func (h *healthStore) MarkDataTypeAvailable(name connector.ExchangeName, dataType healthTypes.DataType, available bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &healthTypes.DataTypeHealth{}
	}

	conn.DataTypes[dataType].Available = available
	conn.LastHealthCheck = time.Now()
}

func (h *healthStore) GetConnectorHealth(name connector.ExchangeName) (*healthTypes.ConnectorHealth, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, exists := h.connectors[name]
	if !exists {
		return nil, false
	}

	connCopy := *conn
	connCopy.DataTypes = make(map[healthTypes.DataType]*healthTypes.DataTypeHealth)
	for k, v := range conn.DataTypes {
		dtCopy := *v
		connCopy.DataTypes[k] = &dtCopy
	}

	return &connCopy, true
}

func (h *healthStore) GetSystemHealth() *healthTypes.SystemHealth {
	h.mu.RLock()
	defer h.mu.RUnlock()

	health := &healthTypes.SystemHealth{
		Connectors:        make(map[connector.ExchangeName]*healthTypes.ConnectorHealth),
		TotalConnectors:   len(h.connectors),
		HealthyConnectors: 0,
		OverallState:      healthTypes.StateConnected,
		StartedAt:         h.startedAt,
	}

	for name, conn := range h.connectors {
		connCopy := *conn
		connCopy.DataTypes = make(map[healthTypes.DataType]*healthTypes.DataTypeHealth)
		for k, v := range conn.DataTypes {
			dtCopy := *v
			connCopy.DataTypes[k] = &dtCopy
		}
		health.Connectors[name] = &connCopy

		if conn.State == healthTypes.StateConnected {
			health.HealthyConnectors++
		}
	}

	if health.HealthyConnectors == 0 {
		health.OverallState = healthTypes.StateDisconnected
	} else if health.HealthyConnectors < health.TotalConnectors {
		health.OverallState = healthTypes.StateDegraded
	}

	return health
}

func (h *healthStore) GetAvailableDataTypes(name connector.ExchangeName) []healthTypes.DataType {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, exists := h.connectors[name]
	if !exists {
		return nil
	}

	var available []healthTypes.DataType
	for dataType, health := range conn.DataTypes {
		if health.Available {
			available = append(available, dataType)
		}
	}

	return available
}

func (h *healthStore) IsDataTypeHealthy(name connector.ExchangeName, dataType healthTypes.DataType) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, exists := h.connectors[name]
	if !exists {
		return false
	}

	dt, dtExists := conn.DataTypes[dataType]
	if !dtExists {
		return false
	}

	return dt.Available && time.Since(dt.LastReceived) < 30*time.Second
}

func (h *healthStore) GetUnhealthyConnectors() []connector.ExchangeName {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var unhealthy []connector.ExchangeName
	for name, conn := range h.connectors {
		if conn.State != healthTypes.StateConnected {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}

func (h *healthStore) GetDegradedDataTypes() map[connector.ExchangeName][]healthTypes.DataType {
	h.mu.RLock()
	defer h.mu.RUnlock()

	degraded := make(map[connector.ExchangeName][]healthTypes.DataType)

	for name, conn := range h.connectors {
		var degradedTypes []healthTypes.DataType
		for dataType, health := range conn.DataTypes {
			if !health.Available || time.Since(health.LastReceived) > 30*time.Second {
				degradedTypes = append(degradedTypes, dataType)
			}
		}
		if len(degradedTypes) > 0 {
			degraded[name] = degradedTypes
		}
	}

	return degraded
}

func (h *healthStore) HasReceivedData(name connector.ExchangeName, dataType healthTypes.DataType) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, exists := h.connectors[name]
	if !exists {
		return false
	}

	dt, dtExists := conn.DataTypes[dataType]
	if !dtExists {
		return false
	}

	return !dt.LastReceived.IsZero()
}

func (h *healthStore) WaitForFirstData(name connector.ExchangeName, dataType healthTypes.DataType, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if h.HasReceivedData(name, dataType) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for %s data from %s", dataType, name)
}
