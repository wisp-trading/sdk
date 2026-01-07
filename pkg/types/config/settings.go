package config

type Configuration interface {
	LoadSettings() (*Settings, error)
	GetConnectors() ([]Connector, error)
	GetEnabledConnectors() ([]Connector, error)

	// Write operations
	SaveSettings(settings *Settings) error
	AddConnector(connector Connector) error
	UpdateConnector(connector Connector) error
	RemoveConnector(name string) error
	EnableConnector(name string, enabled bool) error
}

// Settings represents the main settings structure
type Settings struct {
	Version    string         `mapstructure:"version"`
	Backtest   BacktestConfig `mapstructure:"backtest"`
	Connectors []Connector    `mapstructure:"connectors"`
}

type Connector struct {
	Name        string            `yaml:"name"`
	Enabled     bool              `yaml:"enabled"`
	Network     string            `yaml:"network,omitempty"`
	Assets      []string          `yaml:"assets"`
	Credentials map[string]string `yaml:"credentials"`
}

// BacktestConfig holds backtest settings
type BacktestConfig struct {
	Strategy   string                 `mapstructure:"strategy"`
	Exchange   string                 `mapstructure:"exchange"`
	Pair       string                 `mapstructure:"pair"`
	Timeframe  TimeframeConfig        `mapstructure:"timeframe"`
	Parameters map[string]interface{} `mapstructure:"parameters"`
	Execution  ExecutionConfig        `mapstructure:"execution"`
	Output     OutputConfig           `mapstructure:"output"`
}

// TimeframeConfig defines the backtest time period
type TimeframeConfig struct {
	Start string `mapstructure:"start"`
	End   string `mapstructure:"end"`
}

// ExecutionConfig defines execution parameters
type ExecutionConfig struct {
}

// OutputConfig defines output settings
type OutputConfig struct {
	Format      string `mapstructure:"format"`
	SaveResults bool   `mapstructure:"save_results"`
	ResultsDir  string `mapstructure:"results_dir"`
}

// LiveConfig holds live trading settings
type LiveConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Exchange  string `mapstructure:"exchange"`
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
}
