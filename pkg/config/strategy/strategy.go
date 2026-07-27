package strategy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wisp-trading/sdk/pkg/types/config"
	"gopkg.in/yaml.v3"
)

type strategyConfig struct {
}

func NewStrategyConfigService() config.StrategyConfig {
	return &strategyConfig{}
}

// Load loads and parses a strategy config.yml file
func (c *strategyConfig) Load(path string) (*config.Strategy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Strategy
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize assets map if nil
	if cfg.Assets == nil {
		cfg.Assets = map[string][]config.Asset{}
	}

	// Initialize parameters map if nil
	if cfg.Parameters == nil {
		cfg.Parameters = make(map[string]interface{})
	}

	return &cfg, nil
}

// Save saves a strategy config to config.yml
func (c *strategyConfig) Save(path string, config *config.Strategy) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindStrategies scans the ./strategies directory for available strategies.
// Quiet: no stdout (TUI hosts). Empty directory returns an empty slice, not an error.
func (c *strategyConfig) FindStrategies() ([]config.Strategy, error) {
	strategiesDir := "./strategies"

	if _, err := os.Stat(strategiesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("strategies directory not found: %s (run wisp from a project root, or: wisp init my-bot)", strategiesDir)
	}

	entries, err := os.ReadDir(strategiesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read strategies directory: %w", err)
	}

	var strategies []config.Strategy

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		strategyName := entry.Name()
		strategyPath := filepath.Join(strategiesDir, strategyName)
		configPath := filepath.Join(strategyPath, "config.yml")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		cfg, err := c.Load(configPath)
		if err != nil {
			// Invalid config.yml — skip quietly; user can fix the file and refresh.
			continue
		}

		cfg.Path = strategyPath
		strategies = append(strategies, *cfg)
	}

	// Empty is a valid state (first-run / fresh project) — TUI shows Create CTA.
	if strategies == nil {
		strategies = []config.Strategy{}
	}
	return strategies, nil
}
