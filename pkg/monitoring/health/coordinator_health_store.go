package health

import (
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	health2 "github.com/backtesting-org/kronos-sdk/pkg/types/monitoring/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type coordinatorHealthStore struct {
	connectors        map[connector.ExchangeName]*coordinatorConnectorHealth
	timeProvider      temporal.TimeProvider
	connectorRegistry registry.ConnectorRegistry
	mu                sync.RWMutex
}

type coordinatorConnectorHealth struct {
	DataTypes map[health2.DataType]*health2.DataTypeHealth
}

func NewCoordinatorHealthStore(
	timeProvider temporal.TimeProvider,
	connectorRegistry registry.ConnectorRegistry,
) health2.CoordinatorHealthStore {
	store := &coordinatorHealthStore{
		timeProvider:      timeProvider,
		connectorRegistry: connectorRegistry,
		connectors:        make(map[connector.ExchangeName]*coordinatorConnectorHealth),
	}

	// Initialize from registry - self-contained, no external registration
	connectors := connectorRegistry.GetReadyWebSocketConnectors()
	for _, socketConnector := range connectors {
		name := socketConnector.GetConnectorInfo().Name
		store.connectors[name] = &coordinatorConnectorHealth{
			DataTypes: make(map[health2.DataType]*health2.DataTypeHealth),
		}
	}

	return store
}

func (c *coordinatorHealthStore) RecordDataReceived(name connector.ExchangeName, dataType health2.DataType, source health2.DataSourceType, latency time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, exists := c.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &health2.DataTypeHealth{}
	}

	dt := conn.DataTypes[dataType]
	dt.Available = true
	dt.Source = source
	dt.LastReceived = time.Now()
	dt.Latency = latency
	dt.RecordCount++
	dt.ErrorCount = 0
	dt.LastError = nil
}

func (c *coordinatorHealthStore) RecordDataError(name connector.ExchangeName, dataType health2.DataType, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, exists := c.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &health2.DataTypeHealth{}
	}

	dt := conn.DataTypes[dataType]
	dt.LastError = err
	dt.ErrorCount++

	if dt.ErrorCount >= 3 {
		dt.Available = false
	}
}

func (c *coordinatorHealthStore) MarkDataTypeAvailable(name connector.ExchangeName, dataType health2.DataType, available bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, exists := c.connectors[name]
	if !exists {
		return
	}

	if _, dtExists := conn.DataTypes[dataType]; !dtExists {
		conn.DataTypes[dataType] = &health2.DataTypeHealth{}
	}

	conn.DataTypes[dataType].Available = available
}

func (c *coordinatorHealthStore) GetAvailableDataTypes(name connector.ExchangeName) []health2.DataType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, exists := c.connectors[name]
	if !exists {
		return nil
	}

	var available []health2.DataType
	for dataType, dtHealth := range conn.DataTypes {
		if dtHealth.Available {
			available = append(available, dataType)
		}
	}

	return available
}

func (c *coordinatorHealthStore) IsDataTypeHealthy(name connector.ExchangeName, dataType health2.DataType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, exists := c.connectors[name]
	if !exists {
		return false
	}

	dt, dtExists := conn.DataTypes[dataType]
	if !dtExists {
		return false
	}

	return dt.Available && time.Since(dt.LastReceived) < 30*time.Second
}

func (c *coordinatorHealthStore) HasReceivedData(name connector.ExchangeName, dataType health2.DataType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, exists := c.connectors[name]
	if !exists {
		return false
	}

	dt, dtExists := conn.DataTypes[dataType]
	if !dtExists {
		return false
	}

	return !dt.LastReceived.IsZero()
}

func (c *coordinatorHealthStore) WaitForFirstData(name connector.ExchangeName, dataType health2.DataType, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if c.HasReceivedData(name, dataType) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for %s data from %s", dataType, name)
}

func (c *coordinatorHealthStore) GetConnectorDataHealth(name connector.ExchangeName) map[health2.DataType]*health2.DataTypeHealth {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, exists := c.connectors[name]
	if !exists {
		return nil
	}

	// Return a copy
	result := make(map[health2.DataType]*health2.DataTypeHealth)
	for k, v := range conn.DataTypes {
		dtCopy := *v
		result[k] = &dtCopy
	}

	return result
}

// GetDegradedDataTypes returns data types that have errors (unavailable or with error count)
func (c *coordinatorHealthStore) GetDegradedDataTypes() map[connector.ExchangeName][]health2.DataType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	degraded := make(map[connector.ExchangeName][]health2.DataType)

	for connName, connHealth := range c.connectors {
		for dataType, dtHealth := range connHealth.DataTypes {
			// Degraded = data type is unavailable or has errors
			if !dtHealth.Available || dtHealth.ErrorCount > 0 {
				degraded[connName] = append(degraded[connName], dataType)
			}
		}
	}

	return degraded
}

// GetErrorReport returns aggregated data flow errors
func (c *coordinatorHealthStore) GetErrorReport() *health2.DataFlowErrorReport {
	c.mu.RLock()
	defer c.mu.RUnlock()

	report := &health2.DataFlowErrorReport{
		Errors: make(map[string]map[string]health2.DataFlowError),
	}

	for connName, connHealth := range c.connectors {
		for dataType, dtHealth := range connHealth.DataTypes {
			if !dtHealth.Available || dtHealth.ErrorCount > 0 {
				if _, exists := report.Errors[string(connName)]; !exists {
					report.Errors[string(connName)] = make(map[string]health2.DataFlowError)
				}

				report.Errors[string(connName)][string(dataType)] = health2.DataFlowError{
					Error:      dtHealth.LastError,
					ErrorTime:  dtHealth.LastReceived.Unix(),
					ErrorCount: dtHealth.ErrorCount,
				}
			}
		}
	}

	return report
}
