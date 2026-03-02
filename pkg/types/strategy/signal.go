package strategy

import (
	"time"

	"github.com/google/uuid"
)

// SignalFactory creates market-type-specific signal builders.
type SignalFactory interface {
	// NewSpot creates a builder for spot market signals.
	NewSpot(strategyName StrategyName) SpotSignalBuilder

	// NewPerp creates a builder for perpetual futures signals.
	NewPerp(strategyName StrategyName) PerpSignalBuilder

	// NewPrediction creates a builder for prediction market signals.
	NewPrediction(strategyName StrategyName) PredictionSignalBuilder
}

// Signal is the common interface for all market-type signals.
// The executor uses this to accept any signal type and dispatch accordingly.
type Signal interface {
	GetID() uuid.UUID
	GetStrategy() StrategyName
	GetTimestamp() time.Time
}

// SpotSignal is a trading signal for spot markets.
type SpotSignal struct {
	ID        uuid.UUID     `json:"id"`
	Strategy  StrategyName  `json:"strategy"`
	Timestamp time.Time     `json:"timestamp"`
	Actions   []*SpotAction `json:"actions"`
}

func (s *SpotSignal) GetID() uuid.UUID          { return s.ID }
func (s *SpotSignal) GetStrategy() StrategyName { return s.Strategy }
func (s *SpotSignal) GetTimestamp() time.Time   { return s.Timestamp }

// PerpSignal is a trading signal for perpetual futures markets.
type PerpSignal struct {
	ID        uuid.UUID     `json:"id"`
	Strategy  StrategyName  `json:"strategy"`
	Timestamp time.Time     `json:"timestamp"`
	Actions   []*PerpAction `json:"actions"`
}

func (s *PerpSignal) GetID() uuid.UUID          { return s.ID }
func (s *PerpSignal) GetStrategy() StrategyName { return s.Strategy }
func (s *PerpSignal) GetTimestamp() time.Time   { return s.Timestamp }

// PredictionSignal is a trading signal for prediction markets.
type PredictionSignal struct {
	ID        uuid.UUID           `json:"id"`
	Strategy  StrategyName        `json:"strategy"`
	Timestamp time.Time           `json:"timestamp"`
	Actions   []*PredictionAction `json:"actions"`
}

func (s *PredictionSignal) GetID() uuid.UUID          { return s.ID }
func (s *PredictionSignal) GetStrategy() StrategyName { return s.Strategy }
func (s *PredictionSignal) GetTimestamp() time.Time   { return s.Timestamp }
