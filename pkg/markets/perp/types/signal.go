package types

import (
	"time"

	"github.com/google/uuid"
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// PerpSignal carries perp-domain pair actions (may include leverage).
// Distinct concrete type from SpotSignal so the top-level executor type-switch works.
type PerpSignal interface {
	strategy.Signal
	GetActions() []PerpAction
}

// perpSignal embeds shared PairSignalData.
type perpSignal struct {
	*baseTypes.PairSignalData
}

func (s *perpSignal) GetActions() []PerpAction {
	return s.ActionsCopy()
}

// NewPerpSignal constructs a frozen PerpSignal.
func NewPerpSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []PerpAction) PerpSignal {
	return &perpSignal{PairSignalData: baseTypes.NewPairSignalData(id, name, ts, actions)}
}
