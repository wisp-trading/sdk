package execution

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// HookPlugin is the interface that user's hooks.so must implement
type HookPlugin interface {
	// Name returns the plugin name
	Name() string

	// Version returns the plugin version
	Version() string

	// CreateHooks returns the execution hooks provided by this plugin
	CreateHooks() []ExecutionHook
}

// Executor is the core interface for executing trading signals
type Executor interface {
	// ExecuteSignal processes a signal and executes the associated actions
	ExecuteSignal(signal *strategy.Signal) error

	// HandleTradeExecution is called when a trade is executed to record it
	HandleTradeExecution(trade connector.Trade) error

	// RegisterHook adds an execution hook
	RegisterHook(hook ExecutionHook)
}

// ExecutionHook allows customization of the execution pipeline
type ExecutionHook interface {
	// BeforeExecute is called before executing a signal
	// Returning an error will cancel the execution
	BeforeExecute(ctx *ExecutionContext) error

	// AfterExecute is called after successful execution
	AfterExecute(ctx *ExecutionContext, result *ExecutionResult) error

	// OnError is called when an error occurs during execution
	OnError(ctx *ExecutionContext, err error) error
}
