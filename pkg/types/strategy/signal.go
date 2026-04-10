package strategy

import (
	"time"

	"github.com/google/uuid"
)

// Signal is the common base interface for all market-type signals.
// The executor receives this and type-switches to the concrete domain signal.
type Signal interface {
	GetID() uuid.UUID
	GetStrategy() StrategyName
	GetTimestamp() time.Time
}
