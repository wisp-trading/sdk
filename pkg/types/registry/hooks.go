package registry

import "github.com/wisp-trading/wisp/pkg/types/execution"

// Hooks manages execution hooks at runtime
type Hooks interface {
	// RegisterHook adds a hook to the registry at runtime
	RegisterHook(hook execution.ExecutionHook)

	// RegisterHooks adds multiple hooks to the registry at runtime
	RegisterHooks(hooks []execution.ExecutionHook)

	// GetHooks returns a snapshot of all registered hooks
	GetHooks() []execution.ExecutionHook
}
