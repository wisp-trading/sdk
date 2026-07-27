// Package settings implements config.Configuration: read/write of the global
// connector credentials file (~/.wisp/connectors.yml).
//
// Distinct from:
//   - strategy config.yml (pkg/config/strategy) — exchanges/assets only
//   - StartupConfigLoader — merges both for StartStandalone
package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/wisp-trading/sdk/pkg/types/config"
	"gopkg.in/yaml.v3"
)

// ConfigOptions configures NewConfiguration.
type ConfigOptions struct {
	// SettingsPath is an explicit path to connectors credentials YAML.
	// Empty → ResolveSettingsPath (default ~/.wisp/connectors.yml).
	// Note: StartStandalone(settingsPath) also passes a path into LoadForStrategy;
	// both should agree (CLI leaves this empty and uses StartStandalone's arg).
	SettingsPath string
}

type settings struct {
	settingsPath string // active credentials file path
	cache        *config.Settings
}

// NewConfiguration creates the global connector-credentials store.
func NewConfiguration(opts ConfigOptions) config.Configuration {
	path := config.ResolveSettingsPath(opts.SettingsPath)
	return &settings{
		settingsPath: path,
	}
}

// LoadSettings loads connector credentials from path (or the configured default).
// Reloads when path differs from the cached file.
func (c *settings) LoadSettings(path string) (*config.Settings, error) {
	loadPath := path
	if loadPath == "" {
		loadPath = c.settingsPath
	}
	if loadPath == "" {
		loadPath = config.ResolveSettingsPath("")
	}

	if c.cache != nil && loadPath == c.settingsPath {
		return c.cache, nil
	}

	if !c.fileExists(loadPath) {
		return nil, fmt.Errorf(
			"connector settings not found at %s — add keys via wisp → Settings (~/.wisp/connectors.yml), or pass an explicit path",
			loadPath,
		)
	}

	v := viper.New()
	v.SetConfigFile(loadPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read connector settings: %w", err)
	}

	var loaded config.Settings
	if err := v.Unmarshal(&loaded); err != nil {
		return nil, fmt.Errorf("failed to parse connector settings: %w", err)
	}

	c.cache = &loaded
	c.settingsPath = loadPath
	return c.cache, nil
}

// GetConnectors returns exchange credential rows.
// Missing file is not an error — empty list so CLI Settings can add the first key.
func (c *settings) GetConnectors() ([]config.Connector, error) {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return nil, err
	}
	if c.cache.Connectors == nil {
		return []config.Connector{}, nil
	}
	return c.cache.Connectors, nil
}

// GetEnabledConnectors returns enabled connectors only.
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

// SaveSettings writes the credentials file (0600).
func (c *settings) SaveSettings(s *config.Settings) error {
	dir := filepath.Dir(c.settingsPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal connector settings: %w", err)
	}

	if err := os.WriteFile(c.settingsPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write connector settings: %w", err)
	}

	c.cache = s
	return nil
}

// AddConnector appends a connector row and saves.
func (c *settings) AddConnector(connector config.Connector) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

	for _, existing := range c.cache.Connectors {
		if existing.Name == connector.Name {
			return fmt.Errorf("connector with name '%s' already exists", connector.Name)
		}
	}

	c.cache.Connectors = append(c.cache.Connectors, connector)
	return c.SaveSettings(c.cache)
}

// UpdateConnector replaces an existing connector row and saves.
func (c *settings) UpdateConnector(connector config.Connector) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

	found := false
	for i, existing := range c.cache.Connectors {
		if existing.Name == connector.Name {
			c.cache.Connectors[i] = connector
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", connector.Name)
	}

	return c.SaveSettings(c.cache)
}

// RemoveConnector deletes a connector row by name and saves.
func (c *settings) RemoveConnector(name string) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

	filtered := make([]config.Connector, 0, len(c.cache.Connectors))
	found := false
	for _, connector := range c.cache.Connectors {
		if connector.Name != name {
			filtered = append(filtered, connector)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", name)
	}

	c.cache.Connectors = filtered
	return c.SaveSettings(c.cache)
}

// EnableConnector toggles enabled and saves.
func (c *settings) EnableConnector(name string, enabled bool) error {
	if err := c.ensureLoadedOrEmpty(); err != nil {
		return err
	}

	found := false
	for i, connector := range c.cache.Connectors {
		if connector.Name == name {
			c.cache.Connectors[i].Enabled = enabled
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("connector with name '%s' not found", name)
	}

	return c.SaveSettings(c.cache)
}

// ensureLoadedOrEmpty loads credentials, or starts empty if the file is missing
// (first-time CLI Settings under ~/.wisp).
func (c *settings) ensureLoadedOrEmpty() error {
	if c.cache != nil {
		return nil
	}
	if _, err := c.LoadSettings(""); err != nil {
		if !c.fileExists(c.settingsPath) {
			c.cache = &config.Settings{Connectors: nil}
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
