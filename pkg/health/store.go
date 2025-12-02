package health

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	healthTypes "github.com/backtesting-org/kronos-sdk/pkg/types/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// HealthStore aggregates error reporting from connectors and coordinators
type healthStore struct {
	connectorErrors healthTypes.ConnectorErrorStore
	coordinatorData healthTypes.CoordinatorHealthStore
	timeProvider    temporal.TimeProvider
	startedAt       time.Time
}

// NewHealthStore creates a unified health store
func NewHealthStore(
	timeProvider temporal.TimeProvider,
	connectorErrors healthTypes.ConnectorErrorStore,
	coordinatorData healthTypes.CoordinatorHealthStore,
) healthTypes.HealthStore {
	return &healthStore{
		timeProvider:    timeProvider,
		connectorErrors: connectorErrors,
		coordinatorData: coordinatorData,
		startedAt:       timeProvider.Now(),
	}
}

// GetSystemHealth returns the overall system health state
func (h *healthStore) GetSystemHealth() *healthTypes.SystemHealth {
	return &healthTypes.SystemHealth{
		OverallState: healthTypes.StateConnected,
		StartedAt:    h.startedAt,
	}
}

// GetUnhealthyConnectors returns all connectors with connection errors
func (h *healthStore) GetUnhealthyConnectors() []connector.ExchangeName {
	return h.connectorErrors.GetUnhealthyConnectors()
}

// GetDegradedDataTypes returns data types with data flow errors
func (h *healthStore) GetDegradedDataTypes() map[connector.ExchangeName][]healthTypes.DataType {
	return h.coordinatorData.GetDegradedDataTypes()
}

// GetAvailableDataTypes delegates to coordinatorHealthStore
func (h *healthStore) GetAvailableDataTypes(name connector.ExchangeName) []healthTypes.DataType {
	return h.coordinatorData.GetAvailableDataTypes(name)
}

// IsDataTypeHealthy delegates to coordinatorHealthStore
func (h *healthStore) IsDataTypeHealthy(name connector.ExchangeName, dataType healthTypes.DataType) bool {
	return h.coordinatorData.IsDataTypeHealthy(name, dataType)
}

// HasReceivedData delegates to coordinatorHealthStore
func (h *healthStore) HasReceivedData(name connector.ExchangeName, dataType healthTypes.DataType) bool {
	return h.coordinatorData.HasReceivedData(name, dataType)
}

// WaitForFirstData delegates to coordinatorHealthStore
func (h *healthStore) WaitForFirstData(name connector.ExchangeName, dataType healthTypes.DataType, timeout time.Duration) error {
	return h.coordinatorData.WaitForFirstData(name, dataType, timeout)
}
