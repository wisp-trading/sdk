package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// SpotSignal carries the actions a strategy wants to execute in the spot market.
// GetActions returns a copy — mutating the result does not affect the signal.
type SpotSignal interface {
	strategy.Signal
	GetActions() []SpotAction
}

// spotSignal is the unexported concrete implementation of SpotSignal.
type spotSignal struct {
	id        uuid.UUID
	strat     strategy.StrategyName
	timestamp time.Time
	actions   []SpotAction
}

func (s *spotSignal) GetID() uuid.UUID              { return s.id }
func (s *spotSignal) GetStrategy() strategy.StrategyName { return s.strat }
func (s *spotSignal) GetTimestamp() time.Time        { return s.timestamp }
func (s *spotSignal) GetActions() []SpotAction {
	result := make([]SpotAction, len(s.actions))
	copy(result, s.actions)
	return result
}

// NewSpotSignal constructs a frozen SpotSignal. Called by the builder.
// The actions slice is copied so callers cannot mutate the signal after construction.
func NewSpotSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []SpotAction) SpotSignal {
	copied := make([]SpotAction, len(actions))
	copy(copied, actions)
	return &spotSignal{id: id, strat: name, timestamp: ts, actions: copied}
}
