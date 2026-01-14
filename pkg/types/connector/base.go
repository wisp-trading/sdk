package connector

// Connector handles lifecycle and metadata - ALL connectors implement this
type Connector interface {
	GetConnectorInfo() *Info
	Initialize(config Config) error
	IsInitialized() bool
	NewConfig() Config
	SupportsTradingOperations() bool
	SupportsRealTimeData() bool
}
