package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
)

// Manager handles plugin loading and management
type Manager struct {
	storage       Storage
	logger        Logger
	pluginDir     string
	loadedPlugins map[uuid.UUID]*LoadedPlugin
	mu            sync.RWMutex
}

// Config for plugin manager
type Config struct {
	Storage   Storage
	Logger    Logger
	PluginDir string
}

// NewManager creates a new plugin manager
func NewManager(cfg Config) *Manager {
	return &Manager{
		storage:       cfg.Storage,
		logger:        cfg.Logger,
		pluginDir:     cfg.PluginDir,
		loadedPlugins: make(map[uuid.UUID]*LoadedPlugin),
	}
}

// LoadPlugin loads a plugin from a file path and stores its metadata
func (m *Manager) LoadPlugin(ctx context.Context, pluginPath, createdBy string) (*Metadata, error) {
	m.logger.Info("Loading plugin", "path", pluginPath)

	// Validate file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin file does not exist: %s", pluginPath)
	}

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the NewStrategy symbol
	newStrategySymbol, err := p.Lookup("NewStrategy")
	if err != nil {
		// Try alternative symbol: Strategy variable export
		strategySymbol, err := p.Lookup("Strategy")
		if err != nil {
			return nil, fmt.Errorf("plugin must export NewStrategy function or Strategy variable: %w", err)
		}

		// Type assert to strategy.Strategy
		strat, ok := strategySymbol.(*strategy.Strategy)
		if !ok || strat == nil {
			return nil, fmt.Errorf("Strategy symbol is not of type strategy.Strategy")
		}

		// Extract metadata from strategy instance
		metadata := extractMetadata(*strat)
		metadata.ID = uuid.New()
		metadata.PluginPath = pluginPath
		metadata.CreatedBy = createdBy

		// Store in storage
		if err := m.storage.SavePlugin(ctx, metadata); err != nil {
			return nil, fmt.Errorf("failed to store plugin metadata: %w", err)
		}

		m.logger.Info("Plugin loaded successfully", "id", metadata.ID.String(), "name", metadata.Name)
		return metadata, nil
	}

	// Type assert NewStrategy function
	newStrategyFunc, ok := newStrategySymbol.(func() strategy.Strategy)
	if !ok {
		return nil, fmt.Errorf("NewStrategy must be a function returning strategy.Strategy")
	}

	// Create a temporary instance to extract metadata
	tempStrategy := newStrategyFunc()
	if tempStrategy == nil {
		return nil, fmt.Errorf("NewStrategy() returned nil")
	}

	// Extract metadata
	metadata := extractMetadata(tempStrategy)
	metadata.ID = uuid.New()
	metadata.PluginPath = pluginPath
	metadata.CreatedBy = createdBy

	// Store in storage
	if err := m.storage.SavePlugin(ctx, metadata); err != nil {
		return nil, fmt.Errorf("failed to store plugin metadata: %w", err)
	}

	// Cache the loaded plugin
	m.mu.Lock()
	m.loadedPlugins[metadata.ID] = &LoadedPlugin{
		ID:           metadata.ID,
		Name:         metadata.Name,
		Plugin:       p,
		StrategyFunc: newStrategyFunc,
		Metadata:     metadata,
	}
	m.mu.Unlock()

	m.logger.Info("Plugin loaded successfully", "id", metadata.ID.String(), "name", metadata.Name)
	return metadata, nil
}

// GetLoadedPlugin retrieves a loaded plugin by ID
func (m *Manager) GetLoadedPlugin(ctx context.Context, id uuid.UUID) (*LoadedPlugin, error) {
	m.mu.RLock()
	loaded, exists := m.loadedPlugins[id]
	m.mu.RUnlock()

	if exists {
		return loaded, nil
	}

	// Plugin not in memory, try to load from storage and file
	metadata, err := m.storage.GetPlugin(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("plugin not found in storage: %w", err)
	}

	// Load the plugin file
	p, err := plugin.Open(metadata.PluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file: %w", err)
	}

	// Try to get NewStrategy function
	newStrategySymbol, err := p.Lookup("NewStrategy")
	if err != nil {
		return nil, fmt.Errorf("plugin must export NewStrategy function: %w", err)
	}

	newStrategyFunc, ok := newStrategySymbol.(func() strategy.Strategy)
	if !ok {
		return nil, fmt.Errorf("NewStrategy must be a function returning strategy.Strategy")
	}

	// Cache and return
	loaded = &LoadedPlugin{
		ID:           metadata.ID,
		Name:         metadata.Name,
		Plugin:       p,
		StrategyFunc: newStrategyFunc,
		Metadata:     metadata,
	}

	m.mu.Lock()
	m.loadedPlugins[id] = loaded
	m.mu.Unlock()

	return loaded, nil
}

// InstantiateStrategy creates a new strategy instance from a loaded plugin
func (m *Manager) InstantiateStrategy(ctx context.Context, id uuid.UUID) (strategy.Strategy, error) {
	loaded, err := m.GetLoadedPlugin(ctx, id)
	if err != nil {
		return nil, err
	}

	if loaded.StrategyFunc == nil {
		return nil, fmt.Errorf("plugin does not have a valid NewStrategy function")
	}

	strat := loaded.StrategyFunc()
	if strat == nil {
		return nil, fmt.Errorf("NewStrategy() returned nil")
	}

	return strat, nil
}

// ListPlugins retrieves all plugins from storage
func (m *Manager) ListPlugins(ctx context.Context, limit, offset int) ([]*Metadata, error) {
	return m.storage.ListPlugins(ctx, limit, offset)
}

// GetPluginMetadata retrieves plugin metadata by ID
func (m *Manager) GetPluginMetadata(ctx context.Context, id uuid.UUID) (*Metadata, error) {
	return m.storage.GetPlugin(ctx, id)
}

// DeletePlugin removes a plugin
func (m *Manager) DeletePlugin(ctx context.Context, id uuid.UUID) error {
	// Remove from cache
	m.mu.Lock()
	delete(m.loadedPlugins, id)
	m.mu.Unlock()

	// Delete from storage
	if err := m.storage.DeletePlugin(ctx, id); err != nil {
		return err
	}

	m.logger.Info("Plugin deleted", "id", id.String())
	return nil
}

// SavePluginFile saves an uploaded plugin file to the plugin directory
func (m *Manager) SavePluginFile(fileName string, data []byte) (string, error) {
	// Ensure plugin directory exists
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Generate unique filename
	pluginID := uuid.New()
	ext := filepath.Ext(fileName)
	if ext != ".so" {
		return "", fmt.Errorf("invalid plugin file extension: must be .so")
	}

	uniqueFileName := fmt.Sprintf("%s_%s%s", pluginID.String(), filepath.Base(fileName[:len(fileName)-len(ext)]), ext)
	filePath := filepath.Join(m.pluginDir, uniqueFileName)

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write plugin file: %w", err)
	}

	m.logger.Info("Plugin file saved", "path", filePath)
	return filePath, nil
}
