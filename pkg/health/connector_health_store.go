package health

import (
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type connectorHealthStore struct {
	connectorRegistry registry.ConnectorRegistry
	connectors        map[connector.ExchangeName]*connectorStatus
	timeProvider      temporal.TimeProvider
	mu                sync.RWMutex
}

type connectorStatus struct {
	State     health.ConnectionState
	LastError error
	ErrorTime time.Time
}

func NewConnectorHealthStore(
	timeProvider temporal.TimeProvider,
	connectorRegistry registry.ConnectorRegistry,
) health.ConnectorErrorStore {
	errorStore := &connectorHealthStore{
		timeProvider:      timeProvider,
		connectorRegistry: connectorRegistry,
		connectors:        make(map[connector.ExchangeName]*connectorStatus),
	}

	connectors := connectorRegistry.GetReadyWebSocketConnectors()

	for _, socketConnector := range connectors {
		name := socketConnector.GetConnectorInfo().Name
		errorStore.connectors[name] = &connectorStatus{}
	}

	return errorStore
}

func (c *connectorHealthStore) RecordConnectorError(name connector.ExchangeName, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.connectors[name]; !exists {
		c.connectors[name] = &connectorStatus{
			State: health.StateConnecting,
		}
	}

	c.connectors[name].LastError = err
	c.connectors[name].ErrorTime = c.timeProvider.Now()
	c.connectors[name].State = health.StateDegraded
}

func (c *connectorHealthStore) UpdateConnectionState(name connector.ExchangeName, state health.ConnectionState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.connectors[name]; !exists {
		c.connectors[name] = &connectorStatus{}
	}

	c.connectors[name].State = state
}

func (c *connectorHealthStore) GetConnectorState(name connector.ExchangeName) (health.ConnectionState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status, exists := c.connectors[name]
	if !exists {
		return "", false
	}

	return status.State, true
}

func (c *connectorHealthStore) GetConnectorError(name connector.ExchangeName) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status, exists := c.connectors[name]
	if !exists {
		return nil
	}

	return status.LastError
}

func (c *connectorHealthStore) CountTrackedConnectors() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.connectors)
}

// GetUnhealthyConnectors returns all connectors with connection errors
func (c *connectorHealthStore) GetUnhealthyConnectors() []connector.ExchangeName {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var unhealthy []connector.ExchangeName
	for name, status := range c.connectors {
		if status.State != health.StateConnected {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}
