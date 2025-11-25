package plugin

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
)

// Manager defines the interface for plugin management operations
type Manager interface {
	LoadPlugin(ctx context.Context, pluginPath, createdBy string) (*Metadata, error)
	GetLoadedPlugin(ctx context.Context, id uuid.UUID) (*LoadedPlugin, error)
	InstantiateStrategy(ctx context.Context, id uuid.UUID) (strategy.Strategy, error)

	// LoadHookPlugin loads a hook plugin and returns its metadata
	LoadHookPlugin(ctx context.Context, pluginPath, createdBy string) (*Metadata, error)

	// GetLoadedHookPlugin retrieves a loaded hook plugin by ID
	GetLoadedHookPlugin(ctx context.Context, id uuid.UUID) (*LoadedHookPlugin, error)

	// InstantiateHooks creates hook instances from a loaded hook plugin
	InstantiateHooks(ctx context.Context, id uuid.UUID) ([]execution.ExecutionHook, error)

	// ListPlugins retrieves all plugins (strategies and hooks) from storage
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
