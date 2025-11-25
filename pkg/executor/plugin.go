package executor

import (
	"fmt"
	"plugin"

	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
)

// LoadHookPlugin loads a hook plugin from a .so file
func LoadHookPlugin(path string) (execution.HookPlugin, error) {
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	symPlugin, err := plug.Lookup("HookPlugin")
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export HookPlugin symbol: %w", path, err)
	}

	hookPlugin, ok := symPlugin.(execution.HookPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin %s HookPlugin has invalid type", path)
	}

	return hookPlugin, nil
}
