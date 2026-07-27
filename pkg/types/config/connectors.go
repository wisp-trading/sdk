package config

import "github.com/wisp-trading/sdk/pkg/types/connector"

// ConnectorService bridges Settings rows (config.Connector) to SDK connector.Config
// objects via registry.Connector.NewConfig() + Validate().
//
// Flow for live start:
//
//	GetConnectorConfigsForStrategy(strategy.Exchanges)
//	  → match registry + user Settings
//	  → MapToSDKConfig + Validate
//	  → map[ExchangeName]connector.Config for runtime.Initialize
type ConnectorService interface {
	GetMatchingConnectors() (map[connector.ExchangeName]Connector, error)
	ValidateConnectorConfig(exchangeName connector.ExchangeName, userConnector Connector) error
	// MapToSDKConfig maps YAML credentials into the exchange's NewConfig() type (JSON round-trip).
	MapToSDKConfig(userConnector Connector) (connector.Config, error)
	GetConnectorConfigsForStrategy(exchangeNames []string) (map[connector.ExchangeName]connector.Config, error)
	// GetRequiredCredentialFields discovers form fields from NewConfig() JSON tags (not a second schema).
	GetRequiredCredentialFields(exchangeName string) []string
}
