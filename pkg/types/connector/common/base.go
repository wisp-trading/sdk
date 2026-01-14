package common

import "github.com/backtesting-org/kronos-sdk/pkg/types/connector"

// BaseConnector handles lifecycle and metadata - ALL connectors implement this
type BaseConnector interface {
	GetConnectorInfo() *connector.Info
	Initialize(config connector.Config) error
	IsInitialized() bool
	NewConfig() connector.Config
	SupportsTradingOperations() bool
	SupportsRealTimeData() bool
}
