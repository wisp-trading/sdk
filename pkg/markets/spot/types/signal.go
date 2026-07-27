package types

import (
	"time"

	"github.com/google/uuid"
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// SpotSignal carries spot-domain pair actions.
// Distinct concrete type from PerpSignal so the top-level executor type-switch works.
type SpotSignal interface {
	strategy.Signal
	GetActions() []SpotAction
}

// spotSignal embeds shared PairSignalData.
type spotSignal struct {
	*baseTypes.PairSignalData
}

func (s *spotSignal) GetActions() []SpotAction {
	return s.ActionsCopy()
}

// NewSpotSignal constructs a frozen SpotSignal.
func NewSpotSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []SpotAction) SpotSignal {
	return &spotSignal{PairSignalData: baseTypes.NewPairSignalData(id, name, ts, actions)}
}
