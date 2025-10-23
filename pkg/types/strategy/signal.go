package strategy

import (
	"time"

	"github.com/google/uuid"
)

type Signal struct {
	ID        uuid.UUID     `json:"id"`
	Strategy  StrategyName  `json:"strategy"`
	Actions   []TradeAction `json:"actions"`
	Timestamp time.Time     `json:"timestamp"`
}
