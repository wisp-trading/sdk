package config

// Configuration is the store for global connector credentials
// (~/.wisp/connectors.yml by default; see ResolveSettingsPath).
//
// Naming note: this is *not* strategy config.yml. Strategy exchanges/assets
// live in Strategy / StrategyConfig. This interface only manages API keys and
// enable flags shared by all strategies.
type Configuration interface {
	// LoadSettings loads the credentials file. Empty path → settings path from construction / ResolveSettingsPath.
	LoadSettings(path string) (*Settings, error)
	GetConnectors() ([]Connector, error)
	GetEnabledConnectors() ([]Connector, error)

	SaveSettings(settings *Settings) error
	AddConnector(connector Connector) error
	UpdateConnector(connector Connector) error
	RemoveConnector(name string) error
	EnableConnector(name string, enabled bool) error
}

// Settings is the on-disk shape of ~/.wisp/connectors.yml.
// Only connectors (exchange keys) belong here — never strategy parameters.
type Settings struct {
	Version    string      `yaml:"version,omitempty" mapstructure:"version"`
	Connectors []Connector `yaml:"connectors" mapstructure:"connectors"`
}

// Connector is one exchange row in Settings (user credentials + enable flag).
// Distinct from connector.Config (SDK-typed exchange config after MapToSDKConfig)
// and from registry.Connector (live initialized client).
type Connector struct {
	Name        string            `yaml:"name" mapstructure:"name"`
	Enabled     bool              `yaml:"enabled" mapstructure:"enabled"`
	Network     string            `yaml:"network,omitempty" mapstructure:"network"`
	Assets      []string          `yaml:"assets,omitempty" mapstructure:"assets"`
	Credentials map[string]string `yaml:"credentials" mapstructure:"credentials"`
}
