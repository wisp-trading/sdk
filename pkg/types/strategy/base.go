package strategy

import (
	"sync"
	"time"
)

// BaseStrategy provides common functionality for all strategies
type BaseStrategy struct {
	name         StrategyName
	description  string
	riskLevel    RiskLevel
	strategyType StrategyType
	enabled      bool
	createdAt    time.Time
	updatedAt    time.Time
	mu           sync.RWMutex
}

// NewBaseStrategy creates a new base strategy with common fields
func NewBaseStrategy(
	name StrategyName,
	description string,
	riskLevel RiskLevel,
	strategyType StrategyType,
) *BaseStrategy {
	now := time.Now()
	return &BaseStrategy{
		name:         name,
		description:  description,
		riskLevel:    riskLevel,
		strategyType: strategyType,
		enabled:      false, // Start disabled by default
		createdAt:    now,
		updatedAt:    now,
	}
}

// Enable enables the strategy
func (s *BaseStrategy) Enable() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enabled = true
	s.updatedAt = time.Now()
	return nil
}

// Disable disables the strategy
func (s *BaseStrategy) Disable() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enabled = false
	s.updatedAt = time.Now()
	return nil
}

// IsEnabled returns whether the strategy is enabled
func (s *BaseStrategy) IsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enabled
}

// GetName returns the strategy name
func (s *BaseStrategy) GetName() StrategyName {
	return s.name
}

// GetDescription returns the strategy description
func (s *BaseStrategy) GetDescription() string {
	return s.description
}

// GetRiskLevel returns the strategy risk level
func (s *BaseStrategy) GetRiskLevel() RiskLevel {
	return s.riskLevel
}

// GetStrategyType returns the strategy type
func (s *BaseStrategy) GetStrategyType() StrategyType {
	return s.strategyType
}

// GetCreatedAt returns when the strategy was created
func (s *BaseStrategy) GetCreatedAt() time.Time {
	return s.createdAt
}

// GetUpdatedAt returns when the strategy was last updated
func (s *BaseStrategy) GetUpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}
