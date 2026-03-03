package strategy

import (
	"time"

	"github.com/google/uuid"
)

// SignalFactory creates market-type-specific signal builders.
type SignalFactory interface {
	NewSpot(strategyName StrategyName) SpotSignalBuilder
	NewPerp(strategyName StrategyName) PerpSignalBuilder
}

// Signal is the common interface for all market-type signals.
// The executor accepts this and type-switches to the concrete signal type.
type Signal interface {
	GetID() uuid.UUID
	GetStrategy() StrategyName
	GetTimestamp() time.Time
}

// SpotSignal is the interface for spot market signals.
// Callers get type-safe access to spot actions without a type assertion.
type SpotSignal interface {
	Signal
	GetActions() []*SpotAction
}

// PerpSignal is the interface for perpetual futures signals.
type PerpSignal interface {
	Signal
	GetActions() []*PerpAction
}

// spotSignal is the unexported concrete implementation of SpotSignal.
type spotSignal struct {
	id        uuid.UUID
	strategy  StrategyName
	timestamp time.Time
	actions   []*SpotAction
}

func (s *spotSignal) GetID() uuid.UUID          { return s.id }
func (s *spotSignal) GetStrategy() StrategyName { return s.strategy }
func (s *spotSignal) GetTimestamp() time.Time   { return s.timestamp }
func (s *spotSignal) GetActions() []*SpotAction { return s.actions }

// NewSpotSignal constructs a SpotSignal. Used by the builder.
func NewSpotSignal(id uuid.UUID, name StrategyName, ts time.Time, actions []*SpotAction) SpotSignal {
	return &spotSignal{id: id, strategy: name, timestamp: ts, actions: actions}
}

// perpSignal is the unexported concrete implementation of PerpSignal.
type perpSignal struct {
	id        uuid.UUID
	strategy  StrategyName
	timestamp time.Time
	actions   []*PerpAction
}

func (s *perpSignal) GetID() uuid.UUID          { return s.id }
func (s *perpSignal) GetStrategy() StrategyName { return s.strategy }
func (s *perpSignal) GetTimestamp() time.Time   { return s.timestamp }
func (s *perpSignal) GetActions() []*PerpAction { return s.actions }

// NewPerpSignal constructs a PerpSignal. Used by the builder.
func NewPerpSignal(id uuid.UUID, name StrategyName, ts time.Time, actions []*PerpAction) PerpSignal {
	return &perpSignal{id: id, strategy: name, timestamp: ts, actions: actions}
}
