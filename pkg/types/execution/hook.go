package execution

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
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
	// ExecuteSignal processes a trading signal and executes the associated actions
	ExecuteSignal(signal strategy.Signal) error

	// HandleTradeExecution is called when a trade executes on the exchange
	HandleTradeExecution(trade connector.Trade) error
}

// ExecutionHook defines the interface for execution hooks
type ExecutionHook interface {
	// BeforeExecute is called before an order is placed
	// Return an error to prevent the execution
	BeforeExecute(ctx *ExecutionContext) error

	// AfterExecute is called after an order is successfully placed
	AfterExecute(ctx *ExecutionContext, result *ExecutionResult) error

	// OnError is called when an error occurs during execution
	OnError(ctx *ExecutionContext, err error) error
}
