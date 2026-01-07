package config

import "github.com/backtesting-org/kronos-sdk/pkg/types/connector"

type ConnectorService interface {
	FetchAvailableConnectors() []connector.ExchangeName
	GetAvailableConnectorNames() []string
	GetMatchingConnectors() (map[connector.ExchangeName]Connector, error)
	ValidateConnectorConfig(exchangeName connector.ExchangeName, userConnector Connector) error
	MapToSDKConfig(userConnector Connector) (connector.Config, error)
	GetConnectorConfigsForStrategy(exchangeNames []string) (map[connector.ExchangeName]connector.Config, error)
	GetRequiredCredentialFields(exchangeName string) []string
}

type ConnectorAvailability interface {
	IsAvailable(exchange connector.ExchangeName) bool
	ListAvailable() []connector.ExchangeName
	GetConfigType(exchange connector.ExchangeName) connector.Config
}
