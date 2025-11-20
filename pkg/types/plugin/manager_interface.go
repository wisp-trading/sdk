package plugin

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
)

// Manager defines the interface for plugin management operations
type Manager interface {
	LoadPlugin(ctx context.Context, pluginPath, createdBy string) (*Metadata, error)
	GetLoadedPlugin(ctx context.Context, id uuid.UUID) (*LoadedPlugin, error)
	InstantiateStrategy(ctx context.Context, id uuid.UUID) (strategy.Strategy, error)
	ListPlugins(ctx context.Context, limit, offset int) ([]*Metadata, error)
	GetPluginMetadata(ctx context.Context, id uuid.UUID) (*Metadata, error)
	DeletePlugin(ctx context.Context, id uuid.UUID) error
	SavePluginFile(fileName string, data []byte) (string, error)
}

// Config for plugin manager
type Config struct {
	Storage   Storage
	Logger    logging.ApplicationLogger
	PluginDir string
}
