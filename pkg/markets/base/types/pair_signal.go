package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// PairSignalData is the shared frozen payload for spot/perp signals.
// Domain packages embed this in distinct concrete types so the executor
// type-switch can tell markets apart (identical interfaces would collapse).
type PairSignalData struct {
	ID        uuid.UUID
	Strategy  strategy.StrategyName
	Timestamp time.Time
	Actions   []PairAction
}

// NewPairSignalData copies actions into a frozen payload.
func NewPairSignalData(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []PairAction) *PairSignalData {
	copied := make([]PairAction, len(actions))
	copy(copied, actions)
	return &PairSignalData{ID: id, Strategy: name, Timestamp: ts, Actions: copied}
}

func (s *PairSignalData) GetID() uuid.UUID                   { return s.ID }
func (s *PairSignalData) GetStrategy() strategy.StrategyName { return s.Strategy }
func (s *PairSignalData) GetTimestamp() time.Time            { return s.Timestamp }

// ActionsCopy returns a defensive copy of actions.
func (s *PairSignalData) ActionsCopy() []PairAction {
	result := make([]PairAction, len(s.Actions))
	copy(result, s.Actions)
	return result
}
