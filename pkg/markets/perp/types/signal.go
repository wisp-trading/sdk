package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// PerpSignal carries the actions a strategy wants to execute in the perpetual futures market.
// GetActions returns a copy — mutating the result does not affect the signal.
type PerpSignal interface {
	strategy.Signal
	GetActions() []PerpAction
}

// perpSignal is the unexported concrete implementation of PerpSignal.
type perpSignal struct {
	id        uuid.UUID
	strat     strategy.StrategyName
	timestamp time.Time
	actions   []PerpAction
}

func (s *perpSignal) GetID() uuid.UUID              { return s.id }
func (s *perpSignal) GetStrategy() strategy.StrategyName { return s.strat }
func (s *perpSignal) GetTimestamp() time.Time        { return s.timestamp }
func (s *perpSignal) GetActions() []PerpAction {
	result := make([]PerpAction, len(s.actions))
	copy(result, s.actions)
	return result
}

// NewPerpSignal constructs a frozen PerpSignal. Called by the builder.
// The actions slice is copied so callers cannot mutate the signal after construction.
func NewPerpSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []PerpAction) PerpSignal {
	copied := make([]PerpAction, len(actions))
	copy(copied, actions)
	return &perpSignal{id: id, strat: name, timestamp: ts, actions: copied}
}
