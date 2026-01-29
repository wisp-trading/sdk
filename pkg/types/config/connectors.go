package config

import "github.com/wisp-trading/wisp/pkg/types/connector"

type ConnectorService interface {
	GetMatchingConnectors() (map[connector.ExchangeName]Connector, error)
	ValidateConnectorConfig(exchangeName connector.ExchangeName, userConnector Connector) error
	MapToSDKConfig(userConnector Connector) (connector.Config, error)
	GetConnectorConfigsForStrategy(exchangeNames []string) (map[connector.ExchangeName]connector.Config, error)
	GetRequiredCredentialFields(exchangeName string) []string
}
