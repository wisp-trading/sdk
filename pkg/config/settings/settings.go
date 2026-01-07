package settings

import (
	"fmt"
	"os"

	"github.com/backtesting-org/kronos-sdk/pkg/types/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type settings struct {
	SettingsPath string
	settings     *config.Settings
}

func NewConfiguration() config.Configuration {
	return &settings{
		SettingsPath: config.KronosConfigurationPath + ".yml",
	}
}

// LoadSettings loads the Settings settings from the specified file path
func (c *settings) LoadSettings() (*config.Settings, error) {
	if c.settings != nil {
		return c.settings, nil
	}

	if !c.fileExists(c.SettingsPath) {
		return nil, fmt.Errorf("settings instance not found, please run 'settings init' to create one")
	}

	v := viper.New()
	v.SetConfigFile(c.SettingsPath)
	v.SetConfigType("yaml")

	// Enable environment variable substitution
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var settings config.Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Cache the config
	c.settings = &settings

	return c.settings, nil
}

// GetConnectors returns the cached exchange credentials from settings.yml
// If not loaded yet, it will load the settings config first
func (c *settings) GetConnectors() ([]config.Connector, error) {
	if c.settings != nil {
		return c.settings.Connectors, nil
	}

	// Load the full config which will also cache credentials
	if _, err := c.LoadSettings(); err != nil {
		return nil, err
	}

	return c.settings.Connectors, nil
}

// GetEnabledConnectors returns all enabled connectors
func (c *settings) GetEnabledConnectors() ([]config.Connector, error) {
	if c.settings == nil {
		// Load the full config which will also cache settings
		if _, err := c.LoadSettings(); err != nil {
			return nil, err
		}
	}

	enabled := make([]config.Connector, 0)
	for _, ex := range c.settings.Connectors {
		if ex.Enabled {
			enabled = append(enabled, ex)
		}
	}

	return enabled, nil
}

// SaveSettings writes the settings to the kronos.yml file
func (c *settings) SaveSettings(settings *config.Settings) error {
	// Use gopkg.in/yaml.v3 for better control over formatting
	data, err := marshalYAML(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.SettingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	// Update cache
	c.settings = settings

	return nil
}

// AddConnector adds a new connector to the settings
func (c *settings) AddConnector(connector config.Connector) error {
	// Load current settings
	if c.settings == nil {
		if _, err := c.LoadSettings(); err != nil {
			return err
		}
	}

	// Check for duplicate names
	for _, existing := range c.settings.Connectors {
		if existing.Name == connector.Name {
			return fmt.Errorf("connector with name '%s' already exists", connector.Name)
		}
	}

	// Add connector
	c.settings.Connectors = append(c.settings.Connectors, connector)

	// Save
	return c.SaveSettings(c.settings)
}

// UpdateConnector updates an existing connector
func (c *settings) UpdateConnector(connector config.Connector) error {
	// Load current settings
	if c.settings == nil {
		if _, err := c.LoadSettings(); err != nil {
			return err
		}
	}

	// Find and update connector
	found := false
	for i, existing := range c.settings.Connectors {
		if existing.Name == connector.Name {
			c.settings.Connectors[i] = connector
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", connector.Name)
	}

	// Save
	return c.SaveSettings(c.settings)
}

// RemoveConnector removes a connector by name
func (c *settings) RemoveConnector(name string) error {
	// Load current settings
	if c.settings == nil {
		if _, err := c.LoadSettings(); err != nil {
			return err
		}
	}

	// Filter out the connector
	filtered := make([]config.Connector, 0, len(c.settings.Connectors))
	found := false
	for _, connector := range c.settings.Connectors {
		if connector.Name != name {
			filtered = append(filtered, connector)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", name)
	}

	c.settings.Connectors = filtered

	// Save
	return c.SaveSettings(c.settings)
}

// EnableConnector toggles the enabled state of a connector
func (c *settings) EnableConnector(name string, enabled bool) error {
	// Load current settings
	if c.settings == nil {
		if _, err := c.LoadSettings(); err != nil {
			return err
		}
	}

	// Find and update enabled state
	found := false
	for i, connector := range c.settings.Connectors {
		if connector.Name == name {
			c.settings.Connectors[i].Enabled = enabled
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", name)
	}

	// Save
	return c.SaveSettings(c.settings)
}

// FileExists checks if the config file exists
func (c *settings) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// marshalYAML converts Settings to YAML with proper formatting
func marshalYAML(settings *config.Settings) ([]byte, error) {
	return yaml.Marshal(settings)
}
