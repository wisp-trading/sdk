package config

type StrategyConfig interface {
	Load(path string) (*Strategy, error)
	FindStrategies() ([]Strategy, error)
	Save(path string, config *Strategy) error
}

// Asset is a base/quote pair under an exchange in strategy config.yml.
// Domain routing (spot vs perp vs …) is by the connector's MarketType, not YAML.
type Asset struct {
	Base  string `yaml:"base"`
	Quote string `yaml:"quote"`
}

// StrategyExecutionConfig defines strategy execution timing
type StrategyExecutionConfig struct {
	// Interval defines how frequently the strategy will be executed (e.g., "1m", "5m", "1h")
	// If not set, the global tick interval (50ms) is used
	Interval string `yaml:"interval,omitempty"`
}

// Strategy represents the parsed config.yml for a strategy
type Strategy struct {
	Name       string                 `yaml:"name"`
	Path       string                 `yaml:"-"`
	Error      string                 `yaml:"-"`
	Exchanges  []string               `yaml:"exchanges"`
	Assets     map[string][]Asset     `yaml:"assets"`
	Parameters map[string]interface{} `yaml:"parameters"`
}
