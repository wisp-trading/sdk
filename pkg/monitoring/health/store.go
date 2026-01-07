package health

import (
	"time"

	health2 "github.com/backtesting-org/kronos-sdk/pkg/types/monitoring/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// HealthStore aggregates error reports from both stores into a unified view
type healthStore struct {
	connectorErrors health2.ConnectorErrorStore
	coordinatorData health2.CoordinatorHealthStore
	timeProvider    temporal.TimeProvider
	startedAt       time.Time
}

// NewHealthStore creates a unified health reporter
func NewHealthStore(
	timeProvider temporal.TimeProvider,
	connectorErrors health2.ConnectorErrorStore,
	coordinatorData health2.CoordinatorHealthStore,
) health2.HealthStore {
	return &healthStore{
		timeProvider:    timeProvider,
		connectorErrors: connectorErrors,
		coordinatorData: coordinatorData,
		startedAt:       timeProvider.Now(),
	}
}

// GetSystemHealth returns aggregated health report combining both stores
func (h *healthStore) GetSystemHealth() *health2.SystemHealthReport {
	connReport := h.connectorErrors.GetErrorReport()
	dataReport := h.coordinatorData.GetErrorReport()

	hasErrors := len(connReport.Errors) > 0 || len(dataReport.Errors) > 0
	overallState := health2.StateConnected
	if hasErrors {
		overallState = health2.StateDegraded
	}

	return &health2.SystemHealthReport{
		OverallState:    overallState,
		ConnectorErrors: connReport,
		DataFlowErrors:  dataReport,
		StartedAt:       h.startedAt,
		HasErrors:       hasErrors,
	}
}
