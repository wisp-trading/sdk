package plugin

import (
	"context"

	"github.com/google/uuid"
)

// Storage is the interface that applications must implement for plugin persistence
type Storage interface {
	// SavePlugin stores plugin metadata
	SavePlugin(ctx context.Context, metadata *Metadata) error

	// GetPlugin retrieves plugin metadata by ID
	GetPlugin(ctx context.Context, id uuid.UUID) (*Metadata, error)

	// ListPlugins retrieves all plugins with pagination
	ListPlugins(ctx context.Context, limit, offset int) ([]*Metadata, error)

	// DeletePlugin soft deletes a plugin
	DeletePlugin(ctx context.Context, id uuid.UUID) error
}
