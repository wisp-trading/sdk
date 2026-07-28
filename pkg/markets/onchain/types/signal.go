package types

import (
	"time"

	"github.com/google/uuid"
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// OnchainSignal carries onchain-domain pair swap actions.
type OnchainSignal interface {
	strategy.Signal
	GetActions() []OnchainAction
}

type onchainSignal struct {
	*baseTypes.PairSignalData
}

func (s *onchainSignal) GetActions() []OnchainAction {
	return s.ActionsCopy()
}

// NewOnchainSignal constructs a frozen OnchainSignal.
func NewOnchainSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []OnchainAction) OnchainSignal {
	return &onchainSignal{PairSignalData: baseTypes.NewPairSignalData(id, name, ts, actions)}
}
