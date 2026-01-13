package connectors

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/config"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
)

// ValidationService handles validation of settings and connectors
type ValidationService interface {
	ValidateConnectorName(name string) error
	ValidateUniqueNames(connectors []config.Connector) error
	ValidateSettings(config *config.Settings) error
}

type validationService struct {
	connectorRegistry registry.ConnectorRegistry
}

func NewValidationService(connectorRegistry registry.ConnectorRegistry) ValidationService {
	return &validationService{
		connectorRegistry: connectorRegistry,
	}
}

// ValidateConnectorName checks if the connector name is registered
func (v *validationService) ValidateConnectorName(name string) error {
	_, exists := v.connectorRegistry.GetConnector(connector.ExchangeName(name))
	if !exists {
		return fmt.Errorf("connector '%s' is not registered", name)
	}
	return nil
}

// ValidateUniqueNames ensures no duplicate connector names exist
func (v *validationService) ValidateUniqueNames(connectorList []config.Connector) error {
	seen := make(map[string]bool)

	for _, conn := range connectorList {
		if seen[conn.Name] {
			return fmt.Errorf("duplicate connector name: '%s'", conn.Name)
		}
		seen[conn.Name] = true
	}

	return nil
}

// ValidateSettings performs full validation on a Settings object
// Note: This only validates structure, not credentials
// Credential validation happens when actually connecting to the exchange
func (v *validationService) ValidateSettings(config *config.Settings) error {
	// Validate unique names
	if err := v.ValidateUniqueNames(config.Connectors); err != nil {
		return err
	}

	// Validate each connector
	for _, conn := range config.Connectors {
		// Validate name is registered
		if err := v.ValidateConnectorName(conn.Name); err != nil {
			return err
		}

		// Basic structure validation
		if conn.Name == "" {
			return fmt.Errorf("connector name cannot be empty")
		}

		// Validate credentials are present (not empty)
		if conn.Enabled && len(conn.Credentials) == 0 {
			return fmt.Errorf("connector '%s' is enabled but has no credentials configured", conn.Name)
		}
	}

	return nil
}
