package connectors

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/config"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

type connectorService struct {
	config            config.Configuration
	connectorRegistry registry.ConnectorRegistry
}

func NewConnectorService(
	config config.Configuration,
	connectorRegistry registry.ConnectorRegistry,
) config.ConnectorService {
	return &connectorService{
		config:            config,
		connectorRegistry: connectorRegistry,
	}
}

// GetMatchingConnectors returns user-configured connectors that match registered connectors
func (c *connectorService) GetMatchingConnectors() (map[connector.ExchangeName]config.Connector, error) {
	// Get registered connectors from registry
	registeredConnectors := c.connectorRegistry.GetAllBaseConnectors()

	// Create a lookup map for quick checking
	registeredMap := make(map[connector.ExchangeName]bool)
	for _, conn := range registeredConnectors {
		registeredMap[conn.GetConnectorInfo().Name] = true
	}

	// Get user's configured connectors from the settings service
	userConnectors, err := c.config.GetConnectors()
	if err != nil {
		return nil, err
	}

	// Filter to only return matching connectors as a map
	matchingConnectors := make(map[connector.ExchangeName]config.Connector)
	for _, conn := range userConnectors {
		exchangeName := connector.ExchangeName(conn.Name)
		if registeredMap[exchangeName] {
			matchingConnectors[exchangeName] = conn
		}
	}

	return matchingConnectors, nil
}

// ValidateConnectorConfig validates if a specific exchange has the right configuration loaded
func (c *connectorService) ValidateConnectorConfig(exchangeName connector.ExchangeName, userConnector config.Connector) error {
	ve := &ValidationError{
		Exchange:      string(exchangeName),
		Missing:       []string{},
		InvalidFields: make(map[string]string),
	}

	// Check if the connector is registered
	_, exists := c.connectorRegistry.GetConnector(exchangeName)
	if !exists {
		ve.ExchangeNotFound = true
		return ve
	}

	// Check if the connector is enabled
	if !userConnector.Enabled {
		ve.NotEnabled = true
		return ve
	}

	// Check if the user connector name matches the exchange name
	if userConnector.Name != string(exchangeName) {
		ve.InvalidFields["name"] = fmt.Sprintf("expected '%s', got '%s'", exchangeName, userConnector.Name)
		return ve
	}

	// Map user connector to SDK config
	sdkConfig, err := c.MapToSDKConfig(userConnector)
	if err != nil {
		ve.MappingError = err.Error()
		return ve
	}

	// Validate the SDK config using the SDK's own validation logic
	if err := sdkConfig.Validate(); err != nil {
		ve.SDKValidationErr = err.Error()
		return ve
	}

	return nil
}

// MapToSDKConfig maps a user connector configuration to the appropriate SDK config type
func (c *connectorService) MapToSDKConfig(userConnector config.Connector) (connector.Config, error) {
	exchangeName := connector.ExchangeName(userConnector.Name)

	// Get the connector from registry
	conn, exists := c.connectorRegistry.GetConnector(exchangeName)
	if !exists {
		return nil, fmt.Errorf("connector '%s' not registered", exchangeName)
	}

	// Get config from the connector's info
	cfg := conn.NewConfig()
	if cfg == nil {
		return nil, fmt.Errorf("no config found for exchange '%s'", exchangeName)
	}

	// Create a map to hold all the user's configuration data
	configData := make(map[string]interface{})

	// Copy credentials
	for key, value := range userConnector.Credentials {
		configData[key] = value
	}

	// Add network-related fields if present
	if userConnector.Network != "" {
		configData["network"] = userConnector.Network
		configData["use_testnet"] = userConnector.Network == "testnet"
	}

	// Marshal to JSON and unmarshal into the SDK config type
	jsonData, err := json.Marshal(configData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config data: %w", err)
	}

	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into SDK config: %w", err)
	}

	return cfg, nil
}

// GetConnectorConfigsForStrategy returns validated and mapped SDK configs for the given exchange names
func (c *connectorService) GetConnectorConfigsForStrategy(exchangeNames []string) (map[connector.ExchangeName]connector.Config, error) {
	// Get all matching connectors
	allConnectors, err := c.GetMatchingConnectors()
	if err != nil {
		return nil, fmt.Errorf("failed to get connectors: %w", err)
	}

	// Track detailed problems for each exchange
	validationResults := make(map[string]*ValidationError)

	// Filter to only the exchanges this strategy needs and map to SDK configs
	connectorConfigs := make(map[connector.ExchangeName]connector.Config)

	for _, stratExchangeName := range exchangeNames {
		exchangeName := connector.ExchangeName(stratExchangeName)

		// Check if this exchange is registered
		_, exists := c.connectorRegistry.GetConnector(exchangeName)
		if !exists {
			ve := &ValidationError{
				Exchange:         stratExchangeName,
				ExchangeNotFound: true,
			}
			validationResults[stratExchangeName] = ve
			continue
		}

		// Check if this exchange is configured by user
		userConn, exists := allConnectors[exchangeName]
		if !exists {
			ve := &ValidationError{
				Exchange:         stratExchangeName,
				ExchangeNotFound: true,
			}
			validationResults[stratExchangeName] = ve
			continue
		}

		// Check if enabled
		if !userConn.Enabled {
			ve := &ValidationError{
				Exchange:   stratExchangeName,
				NotEnabled: true,
			}
			validationResults[stratExchangeName] = ve
			continue
		}

		// Validate and map to SDK config
		if err := c.ValidateConnectorConfig(exchangeName, userConn); err != nil {
			var valErr *ValidationError
			if errors.As(err, &valErr) {
				validationResults[stratExchangeName] = valErr
			}
			continue
		}

		sdkConfig, err := c.MapToSDKConfig(userConn)
		if err != nil {
			validationResults[stratExchangeName] = &ValidationError{
				Exchange:     stratExchangeName,
				MappingError: err.Error(),
			}
			continue
		}

		connectorConfigs[exchangeName] = sdkConfig
	}

	// If we have any problems, return a detailed error
	if len(validationResults) > 0 {
		return nil, &StrategyValidationError{
			Strategy:            "",
			ExchangeNames:       exchangeNames,
			SuccessfulExchanges: connectorConfigs,
			ValidationErrors:    validationResults,
		}
	}

	return connectorConfigs, nil
}

// GetRequiredCredentialFields returns the credential field names required by an exchange
func (c *connectorService) GetRequiredCredentialFields(exchangeName string) []string {
	conn, exists := c.connectorRegistry.GetConnector(connector.ExchangeName(exchangeName))
	if !exists {
		return []string{}
	}

	cfg := conn.NewConfig()
	if cfg == nil {
		return []string{}
	}

	configBytes, err := json.Marshal(cfg)
	if err != nil {
		return []string{}
	}

	var fieldMap map[string]interface{}
	if err := json.Unmarshal(configBytes, &fieldMap); err != nil {
		return []string{}
	}

	fields := make([]string, 0, len(fieldMap))
	for key := range fieldMap {
		if key != "network" && key != "use_testnet" {
			fields = append(fields, key)
		}
	}

	return fields
}
