package registry

import "github.com/backtesting-org/kronos-sdk/pkg/types/execution"

type Hooks interface {
	RegisterHook(hook execution.ExecutionHook)
	RegisterHooks(hooks []execution.ExecutionHook)
	GetHooks() []execution.ExecutionHook
}
