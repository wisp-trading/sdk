package strategy

import (
	"time"

	"github.com/google/uuid"
)

// SignalFactory creates signal builders.
type SignalFactory interface {
	New(strategyName StrategyName) SignalBuilder
}

type Signal struct {
	ID        uuid.UUID     `json:"id"`
	Strategy  StrategyName  `json:"strategy"`
	Actions   []TradeAction `json:"actions"`
	Timestamp time.Time     `json:"timestamp"`
}
