package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/wisp-trading/sdk/pkg/types/config"
	"gopkg.in/yaml.v3"
)

// ConfigOptions holds configuration for the settings service
type ConfigOptions struct {
	// SettingsPath is an explicit path to connectors credentials YAML.
	// Empty means resolve via config.ResolveSettingsPath (default ~/.wisp/connectors.yml).
	SettingsPath string
}

type settings struct {
	settingsPath string
	settings     *config.Settings
}

// NewConfiguration creates a new configuration service with the given options
func NewConfiguration(opts ConfigOptions) config.Configuration {
	path := config.ResolveSettingsPath(opts.SettingsPath)
	return &settings{
		settingsPath: path,
	}
}

// LoadSettings loads the settings from the given path, or resolved default if empty
func (c *settings) LoadSettings(path string) (*config.Settings, error) {
	if c.settings != nil && path == "" {
		return c.settings, nil
	}

	loadPath := path
	if loadPath == "" {
		loadPath = c.settingsPath
	}
	if loadPath == "" {
		loadPath = config.ResolveSettingsPath("")
	}

	if !c.fileExists(loadPath) {
		return nil, fmt.Errorf(
			"settings file not found at %s — add connectors via the wisp CLI Settings UI (saved to ~/.wisp/connectors.yml), or pass an explicit path",
			loadPath,
		)
	}

	v := viper.New()
	v.SetConfigFile(loadPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var settings config.Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	c.settings = &settings
	c.settingsPath = loadPath

	return c.settings, nil
}

// GetConnectors returns the cached exchange credentials from settings.
// Missing ~/.wisp/connectors.yml is not an error — returns an empty list
// so the Settings UI can show "Add connector" on first run.
func (c *settings) GetConnectors() ([]config.Connector, error) {
	if c.settings != nil {
		if c.settings.Connectors == nil {
			return []config.Connector{}, nil
		}
		return c.settings.Connectors, nil
	}

	if err := c.ensureLoadedOrEmpty(); err != nil {
		return nil, err
	}
	if c.settings.Connectors == nil {
		return []config.Connector{}, nil
	}
	return c.settings.Connectors, nil
}

// GetEnabledConnectors returns all enabled connectors (empty if none / no file yet).
func (c *settings) GetEnabledConnectors() ([]config.Connector, error) {
	all, err := c.GetConnectors()
	if err != nil {
		return nil, err
	}
	enabled := make([]config.Connector, 0)
	for _, ex := range all {
		if ex.Enabled {
			enabled = append(enabled, ex)
		}
	}
	return enabled, nil
}

// SaveSettings writes the settings file (0600). Ensures parent dir exists.
func (c *settings) SaveSettings(settings *config.Settings) error {
	dir := filepath.Dir(c.settingsPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(c.settingsPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	c.settings = settings
	return nil
}

// AddConnector adds a new connector to the settings
func (c *settings) AddConnector(connector config.Connector) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

	for _, existing := range c.settings.Connectors {
		if existing.Name == connector.Name {
			return fmt.Errorf("connector with name '%s' already exists", connector.Name)
		}
	}

	c.settings.Connectors = append(c.settings.Connectors, connector)
	return c.SaveSettings(c.settings)
}

// UpdateConnector updates an existing connector
func (c *settings) UpdateConnector(connector config.Connector) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

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

	return c.SaveSettings(c.settings)
}

// RemoveConnector removes a connector by name
func (c *settings) RemoveConnector(name string) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

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
	return c.SaveSettings(c.settings)
}

// EnableConnector toggles the enabled state of a connector
func (c *settings) EnableConnector(name string, enabled bool) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

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

	return c.SaveSettings(c.settings)
}

// ensureLoadedOrEmpty loads settings, or starts empty if the file does not exist yet
// (first-time CLI Settings use under ~/.wisp).
func (c *settings) ensureLoadedOrEmpty() error {
	if c.settings != nil {
		return nil
	}
	if _, err := c.LoadSettings(""); err != nil {
		if !c.fileExists(c.settingsPath) {
			c.settings = &config.Settings{Connectors: nil}
			return nil
		}
		return err
	}
	return nil
}

func (c *settings) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
