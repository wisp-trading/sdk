package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// OptionsSignal carries the actions a strategy wants to execute in the options market.
// GetActions returns a copy — mutating the result does not affect the signal.
type OptionsSignal interface {
	strategy.Signal
	GetActions() []OptionsAction
}

// optionsSignal is the unexported concrete implementation of OptionsSignal.
type optionsSignal struct {
	id        uuid.UUID
	strat     strategy.StrategyName
	timestamp time.Time
	actions   []OptionsAction
}

func (s *optionsSignal) GetID() uuid.UUID              { return s.id }
func (s *optionsSignal) GetStrategy() strategy.StrategyName { return s.strat }
func (s *optionsSignal) GetTimestamp() time.Time        { return s.timestamp }
func (s *optionsSignal) GetActions() []OptionsAction {
	result := make([]OptionsAction, len(s.actions))
	copy(result, s.actions)
	return result
}

// NewOptionsSignal constructs a frozen OptionsSignal. Called by the builder.
// The actions slice is copied so callers cannot mutate the signal after construction.
func NewOptionsSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []OptionsAction) OptionsSignal {
	copied := make([]OptionsAction, len(actions))
	copy(copied, actions)
	return &optionsSignal{id: id, strat: name, timestamp: ts, actions: copied}
}
