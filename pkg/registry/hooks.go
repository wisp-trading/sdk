package registry

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// hookRegistry maintains a runtime-updateable list of execution hooks
type hookRegistry struct {
	hooks []execution.ExecutionHook
	mu    sync.RWMutex
}

// NewHookRegistry creates a new hook registry
func NewHookRegistry() registry.Hooks {
	return &hookRegistry{
		hooks: make([]execution.ExecutionHook, 0),
	}
}

// RegisterHook adds a hook to the registry at runtime
func (r *hookRegistry) RegisterHook(hook execution.ExecutionHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks = append(r.hooks, hook)
}

// RegisterHooks adds multiple hooks to the registry at runtime
func (r *hookRegistry) RegisterHooks(hooks []execution.ExecutionHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks = append(r.hooks, hooks...)
}

// GetHooks returns a snapshot of all registered hooks
func (r *hookRegistry) GetHooks() []execution.ExecutionHook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make([]execution.ExecutionHook, len(r.hooks))
	copy(snapshot, r.hooks)
	return snapshot
}
