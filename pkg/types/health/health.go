package health

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
)

// DataType represents types of market data
type DataType string

const (
	DataTypeKlines       DataType = "klines"
	DataTypeOrderbooks   DataType = "orderbooks"
	DataTypeTrades       DataType = "trades"
	DataTypeFundingRates DataType = "funding_rates"
	DataTypePositions    DataType = "positions"
)

// DataSourceType indicates how data is being fetched
type DataSourceType string

const (
	SourceWebSocket DataSourceType = "websocket"
	SourceBatch     DataSourceType = "batch"
	SourceHTTP      DataSourceType = "http"
)

// DataTypeHealth tracks health for a specific data type
type DataTypeHealth struct {
	Available    bool
	Source       DataSourceType
	LastReceived time.Time
	LastError    error
	ErrorCount   int
	Latency      time.Duration
	RecordCount  int64
}

// ConnectorHealth represents the health status of a single connector
type ConnectorHealth struct {
	Name            connector.ExchangeName
	State           ConnectionState
	DataTypes       map[DataType]*DataTypeHealth
	LastHealthCheck time.Time
	UptimeSeconds   int64
	ErrorRate       float64
}

// SystemHealth represents overall system health
type SystemHealth struct {
	Connectors        map[connector.ExchangeName]*ConnectorHealth
	TotalConnectors   int
	HealthyConnectors int
	OverallState      ConnectionState
	StartedAt         time.Time
}

// HealthStore manages health status for all connectors and data flows.
// It is a facade that surfaces errors from ConnectorErrorStore and CoordinatorHealthStore.
type HealthStore interface {
	// System health
	GetSystemHealth() *SystemHealth
	GetUnhealthyConnectors() []connector.ExchangeName
	GetDegradedDataTypes() map[connector.ExchangeName][]DataType

	// Data availability checks
	GetAvailableDataTypes(name connector.ExchangeName) []DataType
	IsDataTypeHealthy(name connector.ExchangeName, dataType DataType) bool
	HasReceivedData(name connector.ExchangeName, dataType DataType) bool
	WaitForFirstData(name connector.ExchangeName, dataType DataType, timeout time.Duration) error
}
