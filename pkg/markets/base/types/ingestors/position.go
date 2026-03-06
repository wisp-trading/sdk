package ingestors

import (
	"context"
)

// PositionCoordinator handles trade backfill on startup
type PositionCoordinator interface {
	Start(ctx context.Context) error
	Stop() error
	IsActive() bool
	GetStatus() map[string]interface{}
}
